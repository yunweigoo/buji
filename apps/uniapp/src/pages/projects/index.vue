<script setup lang="ts">
import { onMounted } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()
onMounted(() => workspace.loadWorkspace().catch(() => uni.showToast({ title: '作品加载失败', icon: 'none' })))

function addProject() {
  uni.showModal({
    title: '新建手作计划', editable: true, placeholderText: '例如：秋日托特包',
    success: async ({ confirm, content }) => {
      const name = content?.trim()
      if (!confirm || !name) return
      try { await workspace.createProject(name); uni.showToast({ title: '计划已创建', icon: 'success' }) }
      catch { uni.showToast({ title: '保存失败', icon: 'none' }) }
    },
  })
}
</script>

<template>
  <view class="page">
    <text class="eyebrow">作品集</text><text class="title">每一件完成的手作，都值得被记下。</text><text class="copy">从搭配灵感到完成照片，让布料和作品彼此找到。</text>
    <button class="add" @click="addProject">＋ 新建手作计划</button>
    <view v-if="workspace.projects.length" class="list">
      <view v-for="project in workspace.projects" :key="project.id" class="card"><view><text class="card-title">{{ project.name }}</text><text class="meta">{{ project.kind || '未分类' }}</text></view><text class="status">{{ project.status === 'completed' ? '已完成' : project.status === 'making' ? '制作中' : '计划中' }}</text></view>
    </view>
    <view v-else class="empty">还没有手作计划，先记下一个想做的作品吧。</view>
  </view>
</template>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 42rpx 32rpx 120rpx; background: $paper; }.eyebrow,.title,.copy,.card-title,.meta { display:block; }.eyebrow { color:$sage; font-size:20rpx; font-weight:700; letter-spacing:2rpx; }.title { max-width:620rpx; margin-top:10rpx; color:$ink; font-size:42rpx; font-weight:700; line-height:1.35; }.copy { margin-top:12rpx; color:$muted; font-size:21rpx; }.add { margin:30rpx 0 0; padding:0 26rpx; border-radius:18rpx; background:$sage; color:white; font-size:21rpx; line-height:68rpx; }.list { display:grid; gap:16rpx; margin-top:34rpx; }.card { display:flex; align-items:center; justify-content:space-between; gap:20rpx; padding:26rpx; border:1rpx solid $line; border-radius:26rpx; background:white; }.card-title { color:$ink; font-size:24rpx; font-weight:700; }.meta { margin-top:8rpx; color:$muted; font-size:19rpx; }.status { padding:10rpx 16rpx; border-radius:999rpx; background:$sage-soft; color:$sage-dark; font-size:18rpx; }.empty { margin-top:80rpx; color:$muted; font-size:21rpx; text-align:center; }
</style>
