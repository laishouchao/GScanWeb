<template>
  <div class="task-detail-container">
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header class="header">
        <div class="header-left">
          <h1>内容安全扫描检测服务</h1>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              {{ userStore.user?.username }} <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-if="userStore.isAdmin" @click="navigateToUsers">用户管理</el-dropdown-item>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主体内容 -->
      <el-container>
        <!-- 左侧菜单 -->
        <el-aside width="200px" class="aside">
          <el-menu :default-active="activeMenu" class="menu" router>
            <el-menu-item index="/">
              <el-icon><house /></el-icon>
              <span>首页</span>
            </el-menu-item>
            <el-menu-item index="/tasks">
              <el-icon><menu /></el-icon>
              <span>任务管理</span>
            </el-menu-item>
            <el-menu-item v-if="userStore.isAdmin" index="/users">
              <el-icon><user-filled /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>

        <!-- 右侧内容 -->
        <el-main class="main">
          <!-- 返回按钮 -->
          <el-button @click="navigateToTasks" style="margin-bottom: 20px;">
            <el-icon><arrow-left /></el-icon>
            返回任务列表
          </el-button>

          <!-- 任务详情 -->
          <el-card class="task-detail-card" v-loading="taskStore.loading">
            <template #header>
              <div class="task-detail-header">
                <h2>任务详情</h2>
                <el-button type="primary" @click="exportTask(taskId)">
                  <el-icon><download /></el-icon>
                  导出结果
                </el-button>
              </div>
            </template>
            <div v-if="taskStore.currentTask" class="task-detail-content">
              <el-divider content-position="left">扫描结果</el-divider>
              <el-table :data="taskStore.currentTask" style="width: 100%">
                <el-table-column prop="type" label="类型" width="120" />
                <el-table-column prop="data" label="内容" />
                <el-table-column prop="created_at" label="时间" width="180" />
              </el-table>
            </div>
            <div v-else class="empty-result">
              <el-empty description="暂无扫描结果" />
            </div>
            <el-alert v-if="taskStore.error" :title="taskStore.error" type="error" show-icon :closable="false" style="margin-top: 20px;" />
          </el-card>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useTaskStore } from '../stores/task'
import { ArrowDown, House, Menu, UserFilled, ArrowLeft, Download } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const taskStore = useTaskStore()

const activeMenu = computed(() => route.path)
const taskId = computed(() => parseInt(route.params.id as string))

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const navigateToUsers = () => {
  router.push('/users')
}

const navigateToTasks = () => {
  router.push('/tasks')
}

const exportTask = async (id: number) => {
  try {
    const result = await taskStore.exportTask(id)
    // 下载文件
    if (result.success && result.file) {
      window.open(result.file, '_blank')
      ElMessage.success('导出成功')
    }
  } catch (error) {
    console.error('导出任务失败:', error)
  }
}

onMounted(async () => {
  // 初始化用户信息
  userStore.initialize()
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }
  
  // 加载任务详情
  await taskStore.getTaskDetail(taskId.value)
})
</script>

<style scoped>
.task-detail-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background-color: #409EFF;
  color: white;
  height: 60px;
}

.header-left h1 {
  margin: 0;
  font-size: 20px;
}

.user-info {
  color: white;
  cursor: pointer;
}

.aside {
  background-color: #f5f7fa;
  border-right: 1px solid #e4e7ed;
}

.menu {
  height: 100%;
}

.main {
  padding: 20px;
  background-color: #f0f2f5;
  overflow-y: auto;
}

.task-detail-card {
  margin-bottom: 20px;
}

.task-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-detail-header h2 {
  margin: 0;
  color: #303133;
}

.task-detail-content {
  margin-top: 20px;
}

.empty-result {
  padding: 40px 0;
  text-align: center;
}
</style>
