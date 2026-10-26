<script setup lang="ts">
defineProps<{
  eyebrow: string
  title: string
  description: string
}>()
</script>

<template>
  <view class="page">
    <text class="eyebrow">{{ eyebrow }}</text>
    <text class="title">{{ title }}</text>
    <text class="description">{{ description }}</text>
    <view class="empty-card">
      <text class="empty-icon">✿</text>
      <text class="empty-copy">页面骨架已建立，下一步接入 CloudBase 数据。</text>
    </view>
  </view>
</template>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 52rpx 32rpx 120rpx; background: $paper; }
.eyebrow { display: block; color: $sage; font-size: 20rpx; font-weight: 700; letter-spacing: 2rpx; }
.title { display: block; margin-top: 12rpx; color: $ink; font-size: 44rpx; font-weight: 700; }
.description { display: block; margin-top: 14rpx; color: $muted; font-size: 24rpx; line-height: 1.7; }
.empty-card { min-height: 320rpx; margin-top: 48rpx; display: flex; flex-direction: column; align-items: center; justify-content: center; border: 2rpx dashed #dfe4d9; border-radius: 28rpx; background: #fcfcf9; }
.empty-icon { color: $sage; font-size: 54rpx; }
.empty-copy { margin-top: 18rpx; color: $muted; font-size: 22rpx; }
</style>

