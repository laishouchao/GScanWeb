package api

import (
	"GScan/infoscan/config"
	"GScan/infoscan/dao"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// HttpApi HTTP API服务
type HttpApi struct {
	api     *Api
	userDAO dao.IUserDAO
	config  *config.Config
}

// NewHttpApi 创建HTTP API服务
func NewHttpApi(api *Api, userDAO dao.IUserDAO, config *config.Config) *HttpApi {
	return &HttpApi{
		api:     api,
		userDAO: userDAO,
		config:  config,
	}
}

// Start 启动HTTP API服务
func (h *HttpApi) Start(addr string) error {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 静态文件服务
	r.Static("/static", "./dist")

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.login)
			auth.POST("/register", h.register)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(h.authMiddleware())
		{
			// 任务管理
			tasks := protected.Group("/tasks")
			{
				tasks.POST("/create", h.createTask)
				tasks.GET("/list", h.listTasks)
				tasks.GET("/detail/:id", h.getTaskDetail)
				tasks.POST("/export/:id", h.exportTask)
			}

			// 用户管理
			users := protected.Group("/users")
			users.Use(h.adminMiddleware())
			{
				users.GET("/list", h.listUsers)
				users.POST("/create", h.createUser)
				users.PUT("/update/:id", h.updateUser)
				users.DELETE("/delete/:id", h.deleteUser)
			}
		}
	}

	// 前端路由
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	// 启动服务
	log.Printf("HTTP API服务启动在 %s", addr)
	return r.Run(addr)
}

// 认证中间件
func (h *HttpApi) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单的认证实现，生产环境应使用JWT
		username := c.GetHeader("Username")
		password := c.GetHeader("Password")

		if username == "" || password == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		user, err := h.userDAO.GetUserByUsername(username)
		if err != nil || user.Password != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// 管理员中间件
func (h *HttpApi) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		if user.(*dao.User).Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 登录
func (h *HttpApi) login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	user, err := h.userDAO.GetUserByUsername(loginData.Username)
	if err != nil || user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// 注册
func (h *HttpApi) register(c *gin.Context) {
	var registerData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := h.userDAO.GetUserByUsername(registerData.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建新用户
	newUser := &dao.User{
		Username: registerData.Username,
		Password: registerData.Password, // 生产环境应使用加密密码
		Role:     "user",
	}

	if err := h.userDAO.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功"})
}

// 创建任务
func (h *HttpApi) createTask(c *gin.Context) {
	var taskData struct {
		URLs []string `json:"urls"`
	}

	if err := c.ShouldBindJSON(&taskData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if len(taskData.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL列表不能为空"})
		return
	}

	// 创建任务
	jobName, jobID := h.api.StartCrawlerJob(taskData.URLs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"task": gin.H{
			"id":   jobID,
			"name": jobName,
		},
	})
}

// 列出任务
func (h *HttpApi) listTasks(c *gin.Context) {
	jobs := h.api.GetJobs()
	c.JSON(http.StatusOK, gin.H{"success": true, "tasks": jobs})
}

// 获取任务详情
func (h *HttpApi) getTaskDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	// 获取任务结果
	results := h.api.GetResults(uint(id))
	c.JSON(http.StatusOK, gin.H{"success": true, "result": results})
}

// 导出任务
func (h *HttpApi) exportTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	filename := fmt.Sprintf("task_%d_%s.xlsx", id, time.Now().Format("20060102150405"))
	filepath := filepath.Join(h.config.ResultPath, filename)

	h.api.Out2Excel(uint(id), filepath)

	c.JSON(http.StatusOK, gin.H{"success": true, "file": "/static/" + filename})
}

// 列出用户
func (h *HttpApi) listUsers(c *gin.Context) {
	users, err := h.userDAO.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "users": users})
}

// 创建用户
func (h *HttpApi) createUser(c *gin.Context) {
	var userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	existingUser, _ := h.userDAO.GetUserByUsername(userData.Username)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 创建新用户
	newUser := &dao.User{
		Username: userData.Username,
		Password: userData.Password, // 生产环境应使用加密密码
		Role:     userData.Role,
	}

	if err := h.userDAO.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建用户成功"})
}

// 更新用户
func (h *HttpApi) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var userData struct {
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取用户
	user, err := h.userDAO.GetUserByUsername(c.GetHeader("Username"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新用户信息
	if userData.Password != "" {
		user.Password = userData.Password
	}
	if userData.Role != "" {
		user.Role = userData.Role
	}

	if err := h.userDAO.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新用户成功"})
}

// 删除用户
func (h *HttpApi) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.userDAO.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除用户成功"})
}
