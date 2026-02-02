<template>
  <div class="task-list-container">
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
          <!-- 创建任务按钮 -->
          <el-button type="primary" @click="showCreateTaskDialog = true" style="margin-bottom: 20px;">
            <el-icon><plus /></el-icon>
            创建新任务
          </el-button>

          <!-- 任务列表 -->
          <el-card class="task-list-card">
            <template #header>
              <div class="task-list-header">
                <h2>任务列表</h2>
              </div>
            </template>
            <el-table :data="taskStore.tasks" style="width: 100%" v-loading="taskStore.loading">
              <el-table-column prop="id" label="任务ID" width="80" />
              <el-table-column prop="name" label="任务名称" />
              <el-table-column prop="status" label="状态" width="120">
                <template #default="scope">
                  <el-tag :type="scope.row.status === 'completed' ? 'success' : 'info'">
                    {{ scope.row.status || '未知' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="180" />
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <el-button size="small" @click="navigateToTaskDetail(scope.row.id)">
                    <el-icon><view /></el-icon>
                    查看详情
                  </el-button>
                  <el-button size="small" type="primary" @click="exportTask(scope.row.id)">
                    <el-icon><download /></el-icon>
                    导出结果
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-alert v-if="taskStore.error" :title="taskStore.error" type="error" show-icon :closable="false" style="margin-top: 20px;" />
          </el-card>
        </el-main>
      </el-container>
    </el-container>

    <!-- 创建任务对话框 -->
    <el-dialog v-model="showCreateTaskDialog" title="创建新任务" width="500px">
      <el-form :model="createTaskForm" :rules="createTaskRules" ref="createTaskFormRef">
        <el-form-item label="URL列表" prop="urls">
          <el-input
            v-model="createTaskForm.urls"
            type="textarea"
            :rows="5"
            placeholder="请输入URL列表，每行一个"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateTaskDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateTask" :loading="taskStore.loading">创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useTaskStore } from '../stores/task'
import { ArrowDown, House, Menu, UserFilled, Plus, View, Download } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const taskStore = useTaskStore()

const activeMenu = computed(() => route.path)
const showCreateTaskDialog = ref(false)
const createTaskFormRef = ref<FormInstance>()

const createTaskForm = ref({
  urls: ''
})

const createTaskRules = reactive<FormRules>({
  urls: [
    { required: true, message: '请输入URL列表', trigger: 'blur' }
  ]
})

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

const navigateToUsers = () => {
  router.push('/users')
}

const navigateToTaskDetail = (taskId: number) => {
  router.push(`/task/${taskId}`)
}

const handleCreateTask = async () => {
  if (!createTaskFormRef.value) return
  
  await createTaskFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 解析URL列表
        const urls = createTaskForm.value.urls.split('\n').filter(url => url.trim() !== '')
        await taskStore.createTask(urls)
        showCreateTaskDialog.value = false
        // 重新加载任务列表
        await taskStore.listTasks()
        ElMessage.success('任务创建成功')
      } catch (error) {
        console.error('创建任务失败:', error)
      }
    }
  })
}

const exportTask = async (taskId: number) => {
  try {
    const result = await taskStore.exportTask(taskId)
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
  
  // 加载任务列表
  await taskStore.listTasks()
})
</script>

<style scoped>
.task-list-container {
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

.task-list-card {
  margin-bottom: 20px;
}

.task-list-header h2 {
  margin: 0;
  color: #303133;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>
