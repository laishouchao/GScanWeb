<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="register-header">
          <h2>内容安全扫描检测服务</h2>
          <p>注册账号</p>
        </div>
      </template>
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请确认密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="userStore.loading">注册</el-button>
          <el-button @click="navigateToLogin">返回登录</el-button>
        </el-form-item>
        <el-alert v-if="userStore.error" :title="userStore.error" type="error" show-icon :closable="false" />
        <el-alert v-if="successMessage" :title="successMessage" type="success" show-icon :closable="false" />
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const registerFormRef = ref<FormInstance>()
const successMessage = ref('')

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const registerRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await userStore.register(registerForm.username, registerForm.password)
        successMessage.value = '注册成功，请登录'
        setTimeout(() => {
          router.push('/login')
        }, 2000)
      } catch (error) {
        console.error('注册失败:', error)
      }
    }
  })
}

const navigateToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.register-card {
  width: 400px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
}

.register-header h2 {
  margin-bottom: 10px;
  color: #303133;
}

.register-header p {
  margin: 0;
  color: #606266;
}

.el-form-item {
  margin-bottom: 20px;
}

.el-alert {
  margin-top: 20px;
}
</style>
