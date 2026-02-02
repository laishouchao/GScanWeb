import { defineStore } from 'pinia'
import axios from 'axios'

interface User {
  id: number
  username: string
  role: string
}

interface UserState {
  user: User | null
  token: string | null
  loading: boolean
  error: string | null
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null,
    loading: false,
    error: null
  }),
  getters: {
    isLoggedIn: (state) => !!state.user,
    isAdmin: (state) => state.user?.role === 'admin'
  },
  actions: {
    async login(username: string, password: string) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post('/api/auth/login', {
          username,
          password
        })
        this.user = response.data.user
        this.token = `${username}:${password}` // 简单的认证方式，生产环境应使用JWT
        localStorage.setItem('user', JSON.stringify(this.user))
        localStorage.setItem('token', this.token)
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '登录失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    async register(username: string, password: string) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post('/api/auth/register', {
          username,
          password
        })
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '注册失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('user')
      localStorage.removeItem('token')
    },
    initialize() {
      const userStr = localStorage.getItem('user')
      const tokenStr = localStorage.getItem('token')
      if (userStr && tokenStr) {
        this.user = JSON.parse(userStr)
        this.token = tokenStr
      }
    }
  }
})
