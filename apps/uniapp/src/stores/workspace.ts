import { defineStore } from 'pinia'
import type { Fabric, Project } from '@/types/domain'

interface WorkspaceState {
  fabrics: Fabric[]
  projects: Project[]
  hydrated: boolean
}

export const useWorkspaceStore = defineStore('workspace', {
  state: (): WorkspaceState => ({
    fabrics: [],
    projects: [],
    hydrated: false,
  }),
  getters: {
    fabricCount: (state) => state.fabrics.length,
    completedProjectCount: (state) => state.projects.filter((item) => item.status === 'completed').length,
  },
  actions: {
    hydrateEmptyWorkspace() {
      this.hydrated = true
    },
  },
})

