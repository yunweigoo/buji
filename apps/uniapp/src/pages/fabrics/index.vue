<script setup lang="ts">
import { onMounted } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()

onMounted(async () => {
  try {
    await workspace.loadWorkspace()
  } catch {
    uni.showToast({ title: '布料加载失败', icon: 'none' })
  }
})

function addFabric() {
  uni.showModal({
    title: '记录一款布料',
    editable: true,
    placeholderText: '例如：蓝色碎花棉布',
    success: async ({ confirm, content }) => {
      const name = content?.trim()
      if (!confirm || !name) return
      try {
        await workspace.createFabric(name)
        uni.showToast({ title: '保存成功', icon: 'success' })
      } catch {
        uni.showToast({ title: '保存失败', icon: 'none' })
      }
    },
  })
}
</script>

<template>
  <view class="page">
    <view class="heading">
      <view><text class="eyebrow">我的布料</text><text class="title">我的布料收藏</text><text class="copy">按花色收藏，每一块布都有自己的故事。</text></view>
      <button class="add" @click="addFabric">＋ 记录布料</button>
    </view>
    <view v-if="workspace.loading" class="state">正在翻开布料册…</view>
    <view v-else-if="workspace.error" class="state error">{{ workspace.error }}</view>
    <view v-else-if="workspace.fabrics.length" class="grid">
      <view v-for="fabric in workspace.fabrics" :key="fabric.id" class="card">
        <view class="swatch">✿</view>
        <text class="card-title">{{ fabric.name }}</text>
        <text class="meta">{{ [fabric.material, fabric.pattern].filter(Boolean).join(' · ') || '待补充材质与花纹' }}</text>
        <text v-if="fabric.memo" class="memo">{{ fabric.memo }}</text>
      </view>
    </view>
    <view v-else class="empty"><text class="empty-icon">▦</text><text class="empty-title">布料册还是空的</text><text class="empty-copy">从记录手边第一款布开始吧。</text></view>
  </view>
</template>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 42rpx 32rpx 120rpx; background: $paper; }
.heading { display: flex; align-items: flex-start; justify-content: space-between; gap: 24rpx; }
.eyebrow, .title, .copy, .card-title, .meta, .memo, .empty-title, .empty-copy { display: block; }
.eyebrow { color: $sage; font-size: 20rpx; font-weight: 700; letter-spacing: 2rpx; }
.title { margin-top: 10rpx; color: $ink; font-size: 42rpx; font-weight: 700; }
.copy { margin-top: 10rpx; color: $muted; font-size: 21rpx; }
.add { flex: 0 0 auto; margin: 0; padding: 0 24rpx; border-radius: 18rpx; background: $sage; color: white; font-size: 21rpx; line-height: 66rpx; }
.grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18rpx; margin-top: 36rpx; }
.card { padding: 24rpx; border: 1rpx solid $line; border-radius: 26rpx; background: white; }
.swatch { width: 72rpx; height: 72rpx; display: flex; align-items: center; justify-content: center; border-radius: 22rpx; background: $sage-soft; color: $sage-dark; }
.card-title { margin-top: 20rpx; color: $ink; font-size: 24rpx; font-weight: 700; }
.meta, .memo { margin-top: 8rpx; color: $muted; font-size: 19rpx; line-height: 1.5; }
.state, .empty { margin-top: 80rpx; color: $muted; text-align: center; }.state.error { color: #9d5c50; }
.empty-icon { display: block; color: $sage; font-size: 62rpx; }.empty-title { margin-top: 18rpx; color: $ink; font-size: 26rpx; font-weight: 700; }.empty-copy { margin-top: 10rpx; font-size: 20rpx; }
@media (max-width: 420px) { .heading { display: block; }.add { margin-top: 24rpx; }.grid { grid-template-columns: 1fr; } }
</style>
