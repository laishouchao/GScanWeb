import { defineStore } from 'pinia'
import axios from 'axios'

interface Task {
  id: number
  name: string
  status: string
  created_at: string
}

interface TaskState {
  tasks: Task[]
  currentTask: any | null
  loading: boolean
  error: string | null
}

export const useTaskStore = defineStore('task', {
  state: (): TaskState => ({
    tasks: [],
    currentTask: null,
    loading: false,
    error: null
  }),
  actions: {
    async createTask(urls: string[]) {
      this.loading = true
      this.error = null
      try {
        // 获取认证信息
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        const response = await axios.post('/api/tasks/create', {
          urls
        }, {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '创建任务失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    async listTasks() {
      this.loading = true
      this.error = null
      try {
        // 获取认证信息
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        const response = await axios.get('/api/tasks/list', {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        this.tasks = response.data.tasks
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '获取任务列表失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    async getTaskDetail(id: number) {
      this.loading = true
      this.error = null
      try {
        // 获取认证信息
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        const response = await axios.get(`/api/tasks/detail/${id}`, {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        this.currentTask = response.data.result
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '获取任务详情失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    async exportTask(id: number) {
      this.loading = true
      this.error = null
      try {
        // 获取认证信息
        const user = localStorage.getItem('user')
        if (!user) {
          throw new Error('未认证')
        }
        const userObj = JSON.parse(user)
        const token = localStorage.getItem('token')
        
        const response = await axios.post(`/api/tasks/export/${id}`, {}, {
          headers: {
            Username: userObj.username,
            Password: token?.split(':')[1]
          }
        })
        return response.data
      } catch (error: any) {
        this.error = error.response?.data?.error || '导出任务失败'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
