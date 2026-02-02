<template>
  <div class="home-container">
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
          <el-card class="welcome-card">
            <template #header>
              <div class="welcome-header">
                <h2>欢迎使用内容安全扫描检测服务</h2>
              </div>
            </template>
            <div class="welcome-content">
              <p>本系统基于高性能Go语言内核，提供多用户隔离的内容安全扫描检测功能。</p>
              <el-divider />
              <h3>主要功能</h3>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-card class="feature-card" shadow="hover">
                    <el-icon class="feature-icon"><link /></el-icon>
                    <h4>外部链接检测</h4>
                    <p>扫描网站上的外链，及时发现废弃域名被抢注指向非法网站</p>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card class="feature-card" shadow="hover">
                    <el-icon class="feature-icon"><warning /></el-icon>
                    <h4>敏感信息检测</h4>
                    <p>扫描网站的内容，及时发现敏感信息(如身份证)，避免信息泄露</p>
                  </el-card>
                </el-col>
                <el-col :span="8">
                  <el-card class="feature-card" shadow="hover">
                    <el-icon class="feature-icon"><document /></el-icon>
                    <h4>文件下载检测</h4>
                    <p>扫描网站开放下载的文件连接，对内容进行排查及时发现敏感信息</p>
                  </el-card>
                </el-col>
              </el-row>
              <el-divider />
              <h3>快速开始</h3>
              <el-button type="primary" @click="navigateToTasks">查看任务列表</el-button>
              <el-button @click="navigateToCreateTask">创建新任务</el-button>
            </div>
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
import { ArrowDown, House, Menu, UserFilled, Link, Warning, Document } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

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

const navigateToCreateTask = () => {
  router.push('/tasks')
}

onMounted(() => {
  // 初始化用户信息
  userStore.initialize()
  if (!userStore.isLoggedIn) {
    router.push('/login')
  }
})
</script>

<style scoped>
.home-container {
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

.welcome-card {
  margin-bottom: 20px;
}

.welcome-header h2 {
  margin: 0;
  color: #303133;
}

.welcome-content {
  line-height: 1.6;
}

.feature-card {
  text-align: center;
  padding: 20px;
}

.feature-icon {
  font-size: 48px;
  color: #409EFF;
  margin-bottom: 10px;
}

.feature-card h4 {
  margin: 10px 0;
  color: #303133;
}

.feature-card p {
  color: #606266;
  font-size: 14px;
}
</style>
