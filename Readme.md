# 内容安全扫描检测服务 Web 平台

基于 Go 语言内核的高性能多用户隔离的内容安全扫描检测服务 Web 平台。

## 系统概述

本系统是一个基于 Vue 的高性能多用户隔离的内容安全扫描检测服务 Web 平台，它利用现有的 Go 语言内核代码，提供了一个友好的 Web 界面，使用户可以方便地创建和管理扫描任务，查看扫描结果，并导出报告。

### 核心功能

1. **外部链接检测**：扫描网站上的外链，及时发现废弃域名被抢注指向非法网站
2. **敏感信息检测**：扫描网站的内容，及时发现敏感信息(如身份证)，避免信息泄露
3. **文件下载检测**：扫描网站开放下载的文件连接，对内容进行排查及时发现敏感信息
4. **多用户隔离**：支持多用户登录，每个用户只能查看和管理自己的任务
5. **权限管理**：区分普通用户和管理员权限，管理员可以管理所有用户

## 技术栈

### 后端
- **语言**：Go 1.18+
- **Web 框架**：Gin
- **数据库**：SQLite
- **依赖管理**：Go Modules

### 前端
- **框架**：Vue 3 + TypeScript
- **状态管理**：Pinia
- **路由**：Vue Router
- **UI 组件库**：Element Plus
- **构建工具**：Vite

## 系统架构

系统采用前后端分离的架构设计：

1. **前端**：Vue 3 单页应用，负责用户界面和交互
2. **后端**：Go 语言实现的 RESTful API 服务，负责业务逻辑和数据处理
3. **数据库**：SQLite 数据库，存储用户信息和扫描任务数据

## 快速开始

### 环境要求

- Go 1.18 或更高版本
- Node.js 16 或更高版本
- npm 或 yarn

### 安装和运行

#### 1. 克隆项目

```bash
git clone https://github.com/yourusername/GScanWeb.git
cd GScanWeb
```

#### 2. 安装前端依赖

```bash
npm install
```

#### 3. 构建前端项目

```bash
npm run build
```

#### 4. 运行后端服务

```bash
cd infoscan/cmd
go mod tidy
go build -o infoscan.exe
./infoscan.exe
```

#### 5. 访问系统

打开浏览器，访问 `http://localhost:8080`

### 默认账户

- **用户名**：admin
- **密码**：admin123

## 使用指南

### 1. 登录系统

使用默认账户或注册新账户登录系统。

### 2. 创建扫描任务

1. 在左侧菜单中点击「任务管理」
2. 点击「创建新任务」按钮
3. 在弹出的对话框中输入要扫描的 URL 列表（每行一个）
4. 点击「创建」按钮

### 3. 查看扫描结果

1. 在任务列表中找到刚创建的任务
2. 点击「查看详情」按钮
3. 在任务详情页面查看扫描结果

### 4. 导出扫描报告

1. 在任务列表中找到要导出的任务
2. 点击「导出结果」按钮
3. 系统会生成 Excel 格式的报告并下载

### 5. 用户管理（管理员）

1. 在左侧菜单中点击「用户管理」
2. 可以查看、创建、编辑和删除用户

## 项目结构

```
GScanWeb/
├── infoscan/            # Go 语言内核代码
│   ├── api/            # API 实现
│   ├── cmd/            # 命令行入口
│   ├── config/         # 配置文件
│   ├── dao/            # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── pkg/            # 工具包
├── src/                # Vue 前端代码
│   ├── assets/         # 静态资源
│   ├── components/     # 组件
│   ├── router/         # 路由
│   ├── stores/         # 状态管理
│   ├── views/          # 页面
│   ├── App.vue         # 根组件
│   ├── main.ts         # 入口文件
│   └── style.css       # 全局样式
├── dist/               # 构建输出目录
├── go.mod              # Go 依赖管理
├── package.json        # npm 依赖管理
├── README.md           # 项目文档
└── vite.config.ts      # Vite 配置
```

## 配置说明

### 后端配置

后端配置文件位于 `infoscan/cmd/config.yml`，可以根据需要调整配置：

```yaml
Version: 0.4.10
ResultPath: "./result"      # 结果文件存储路径
SpiderMaxNum: 10           # 最大爬虫数量
Spider:
  Threads: 5               # 爬虫线程数
  Httpspider:
    domain_headers: []
    navigate_timeout_second: 30
    proxy: ""
  page_analyze_timeout_second: 30
  retry: 3
Downloader:
  Enable: false
Name: "InfoScan"
whitelistFile: "whitelist.txt"  # 白名单文件路径
LogPath: "./log"          # 日志文件路径
LogLevel: 0
LogPrintingLevel: 0
```

### 白名单配置

白名单文件位于 `infoscan/cmd/whitelist.txt`，用于过滤安全的外链：

```
edu.cn
gov.cn
```

## 注意事项

1. **性能优化**：对于大型网站的扫描，可能会消耗较多的系统资源，建议根据服务器配置调整爬虫线程数和最大爬虫数量。

2. **安全考虑**：本系统默认使用简单的认证方式，生产环境中建议使用 JWT 等更安全的认证方式，并对密码进行加密存储。

3. **误报处理**：敏感信息检测和关键词检测可能会有一定的误报率，需要人工核查。

4. **数据库备份**：定期备份 SQLite 数据库文件 `data.db`，以防止数据丢失。

5. **日志管理**：系统会生成日志文件，定期清理以避免磁盘空间占用过大。

## 开发指南

### 前端开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 后端开发

```bash
# 安装依赖
cd infoscan/cmd
go mod tidy

# 编译
go build -o infoscan.exe

# 运行
./infoscan.exe
```

## 故障排除

### 常见问题

1. **前端构建失败**
   - 检查 Node.js 版本是否符合要求
   - 尝试删除 `node_modules` 目录并重新安装依赖

2. **后端服务启动失败**
   - 检查 Go 版本是否符合要求
   - 运行 `go mod tidy` 更新依赖
   - 检查端口是否被占用

3. **扫描任务执行失败**
   - 检查目标网站是否可访问
   - 检查网络连接是否正常
   - 查看日志文件了解具体错误信息

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 许可证

MIT License

## 联系方式

- Email: i@vshex.com
- GitHub: https://github.com/Ymjie/GScan
