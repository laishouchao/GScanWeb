<template>
  <div class="user-list-container">
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
            <el-menu-item index="/users">
              <el-icon><user-filled /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>

        <!-- 右侧内容 -->
        <el-main class="main">
          <!-- 创建用户按钮 -->
          <el-button type="primary" @click="showCreateUserDialog = true" style="margin-bottom: 20px;">
            <el-icon><plus /></el-icon>
            创建新用户
          </el-button>

          <!-- 用户列表 -->
          <el-card class="user-list-card">
            <template #header>
              <div class="user-list-header">
                <h2>用户管理</h2>
              </div>
            </template>
            <el-table :data="users" style="width: 100%" v-loading="loading">
              <el-table-column prop="id" label="用户ID" width="80" />
              <el-table-column prop="username" label="用户名" />
              <el-table-column prop="role" label="角色" width="120">
                <template #default="scope">
                  <el-tag :type="scope.row.role === 'admin' ? 'warning' : 'success'">
                    {{ scope.row.role === 'admin' ? '管理员' : '普通用户' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <el-button size="small" @click="handleEditUser(scope.row)">
                    <el-icon><edit /></el-icon>
                    编辑
                  </el-button>
                  <el-button size="small" type="danger" @click="handleDeleteUser(scope.row.id)">
                    <el-icon><delete /></el-icon>
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" style="margin-top: 20px;" />
          </el-card>
        </el-main>
      </el-container>
    </el-container>

    <!-- 创建用户对话框 -->
    <el-dialog v-model="showCreateUserDialog" title="创建新用户" width="500px">
      <el-form :model="userForm" :rules="userRules" ref="userFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="userForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateUserDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateUser" :loading="loading">创建</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog v-model="showEditUserDialog" title="编辑用户" width="500px">
      <el-form :model="editUserForm" :rules="editUserRules" ref="editUserFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="editUserForm.username" placeholder="请输入用户名" disabled />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="editUserForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="editUserForm.role" placeholder="请选择角色">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditUserDialog = false">取消</el-button>
          <el-button type="primary" @click="handleUpdateUser" :loading="loading">更新</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ArrowDown, House, Menu, UserFilled, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const users = ref<any[]>([])
const loading = ref(false)
const error = ref('')

// 创建用户对话框
const showCreateUserDialog = ref(false)
const userFormRef = ref<FormInstance>()
const userForm = ref({
  username: '',
  password: '',
  role: 'user'
})

const userRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'blur' }
  ]
})

// 编辑用户对话框
const showEditUserDialog = ref(false)
const editUserFormRef = ref<FormInstance>()
const editUserForm = ref({
  id: 0,
  username: '',
  password: '',
  role: 'user'
})

const editUserRules = reactive<FormRules>({
  password: [
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'blur' }
  ]
})

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  error.value = ''
  try {
    const user = localStorage.getItem('user')
    if (!user) {
      throw new Error('未认证')
    }
    const userObj = JSON.parse(user)
    const token = localStorage.getItem('token')
    
    const response = await axios.get('/api/users/list', {
      headers: {
        Username: userObj.username,
        Password: token?.split(':')[1]
      }
    })
    users.value = response.data.users
  } catch (err: any) {
    error.value = err.response?.data?.error || '获取用户列表失败'
    console.error('获取用户列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 创建用户
const handleCreateUser = async () => {
  if (!userFormRef.value) return
  
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        await axios.post('/api/users/create', userForm.value, {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        showCreateUserDialog.value = false
        // 重新加载用户列表
        await fetchUsers()
        ElMessage.success('创建用户成功')
      } catch (err: any) {
        error.value = err.response?.data?.error || '创建用户失败'
        console.error('创建用户失败:', err)
      } finally {
        loading.value = false
      }
    }
  })
}

// 编辑用户
const handleEditUser = (user: any) => {
  editUserForm.value = {
    id: user.id,
    username: user.username,
    password: '',
    role: user.role
  }
  showEditUserDialog.value = true
}

// 更新用户
const handleUpdateUser = async () => {
  if (!editUserFormRef.value) return
  
  await editUserFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        await axios.put(`/api/users/update/${editUserForm.value.id}`, {
          password: editUserForm.value.password,
          role: editUserForm.value.role
        }, {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        showEditUserDialog.value = false
        // 重新加载用户列表
        await fetchUsers()
        ElMessage.success('更新用户成功')
      } catch (err: any) {
        error.value = err.response?.data?.error || '更新用户失败'
        console.error('更新用户失败:', err)
      } finally {
        loading.value = false
      }
    }
  })
}

// 删除用户
const handleDeleteUser = async (userId: number) => {
  ElMessageBox.confirm('确定要删除该用户吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value = true
    try {
      const user = localStorage.getItem('user')
      if (!user) {
        throw new Error('未认证')
      }
      const userObj = JSON.parse(user)
      const token = localStorage.getItem('token')
      
      await axios.delete(`/api/users/delete/${userId}`, {
        headers: {
          Username: userObj.username,
          Password: token?.split(':')[1]
        }
      })
      // 重新加载用户列表
      await fetchUsers()
      ElMessage.success('删除用户成功')
    } catch (err: any) {
      error.value = err.response?.data?.error || '删除用户失败'
      console.error('删除用户失败:', err)
    } finally {
      loading.value = false
    }
  }).catch(() => {
    // 取消删除
  })
}

onMounted(async () => {
  // 初始化用户信息
  userStore.initialize()
  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    router.push('/')
    return
  }
  
  // 加载用户列表
  await fetchUsers()
})
</script>

<style scoped>
.user-list-container {
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

.user-list-card {
  margin-bottom: 20px;
}

.user-list-header h2 {
  margin: 0;
  color: #303133;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>
