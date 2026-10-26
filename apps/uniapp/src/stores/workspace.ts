import { defineStore } from 'pinia'
import { apiRequest } from '@/api/client'
import type { Fabric, Project } from '@/types/domain'

interface WorkspaceState {
  fabrics: Fabric[]
  projects: Project[]
  hydrated: boolean
  loading: boolean
  error?: string
}

interface BootstrapData {
  fabrics: Fabric[]
  projects: Project[]
  features: Record<string, boolean>
}

export const useWorkspaceStore = defineStore('workspace', {
  state: (): WorkspaceState => ({
    fabrics: [],
    projects: [],
    hydrated: false,
    loading: false,
    error: undefined,
  }),
  getters: {
    fabricCount: (state) => state.fabrics.length,
    completedProjectCount: (state) => state.projects.filter((item) => item.status === 'completed').length,
  },
  actions: {
    async loadWorkspace(force = false) {
      if (this.loading || (this.hydrated && !force)) return
      this.loading = true
      this.error = undefined
      try {
        const data = await apiRequest<BootstrapData>('/v1/bootstrap')
        this.fabrics = data.fabrics
        this.projects = data.projects
        this.hydrated = true
      } catch (error) {
        this.error = error instanceof Error ? error.message : '数据加载失败'
        throw error
      } finally {
        this.loading = false
      }
    },
    async createFabric(name: string) {
      const item = await apiRequest<Fabric>('/v1/fabrics', { method: 'POST', data: { name } })
      this.fabrics.unshift(item)
      return item
    },
    async createProject(name: string) {
      const item = await apiRequest<Project>('/v1/projects', { method: 'POST', data: { name } })
      this.projects.unshift(item)
      return item
    },
  },
})

