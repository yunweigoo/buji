<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import FabricEntryModal from '@/components/FabricEntryModal.vue'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()
const showEntry = ref(false)
const today = computed(() => new Intl.DateTimeFormat('zh-CN', { month: 'long', day: 'numeric', weekday: 'short' }).format(new Date()))

onMounted(async () => {
  try {
    await workspace.loadWorkspace()
  } catch {
    uni.showToast({ title: '数据库连接失败，请稍后重试', icon: 'none' })
  }
})

function addFabric() {
  showEntry.value = true
}

function go(path: string) {
  uni.switchTab({ url: path })
}
</script>

<template>
  <view class="page">
    <view class="topbar">
      <view>
        <text class="eyebrow">工作台</text>
        <text class="title">把喜欢的布，慢慢变成作品。</text>
        <text class="subtitle">一张布一段记忆，搭配、记录和灵感都在这里。</text>
      </view>
      <text class="date">{{ today }}</text>
    </view>

    <view class="hero">
      <text class="hero-kicker">✿ YOUR LITTLE FABRIC STUDIO</text>
      <text class="hero-title">今天，也给好看的布\n安排一个新去处。</text>
      <text class="hero-copy">先认识手里已有的，再开始下一件喜欢的手作。</text>
      <view class="hero-actions">
        <button class="button primary" @click="addFabric">拍照识别布料</button>
        <button class="button light" @click="go('/pages/match/index')">试试搭配</button>
      </view>
    </view>

    <view class="metrics">
      <view class="metric"><text class="metric-icon">▦</text><view><text class="metric-label">我的花色</text><text class="metric-value">{{ workspace.fabricCount }}<small> 款</small></text></view></view>
      <view class="metric"><text class="metric-icon peach">✂</text><view><text class="metric-label">实物布块</text><text class="metric-value">{{ workspace.pieceCount }}<small> 块</small></text></view></view>
      <view class="metric"><text class="metric-icon pink">♡</text><view><text class="metric-label">完成作品</text><text class="metric-value">{{ workspace.completedProjectCount }}<small> 件</small></text></view></view>
    </view>

    <view class="section-head">
      <view><text class="section-title">常用小工具</text><text class="section-copy">今天想做点什么？</text></view>
    </view>
    <view class="tools">
      <button class="tool" @click="addFabric"><text class="tool-icon peach">📷</text><view><text class="tool-title">拍照收进布料库</text><text class="tool-copy">识别颜色与花纹，确认后保存</text></view><text class="arrow">›</text></button>
      <button class="tool" @click="go('/pages/check/index')"><text class="tool-icon pink">⌕</text><view><text class="tool-title">买之前先看看</text><text class="tool-copy">和自己的布料收藏比一比</text></view><text class="arrow">›</text></button>
      <button class="tool" @click="go('/pages/match/index')"><text class="tool-icon">✧</text><view><text class="tool-title">主布怎么搭？</text><text class="tool-copy">拿手边库存试出新组合</text></view><text class="arrow">›</text></button>
    </view>
    <FabricEntryModal :open="showEntry" @close="showEntry = false" />
  </view>
</template>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 42rpx 32rpx 120rpx; background: $paper; }
.topbar { display: flex; justify-content: space-between; align-items: flex-start; gap: 20rpx; }
.eyebrow { display: block; color: $sage; font-size: 20rpx; font-weight: 700; letter-spacing: 2rpx; }
.title { display: block; margin-top: 10rpx; color: $ink; font-size: 44rpx; font-weight: 700; line-height: 1.3; }
.subtitle { display: block; margin-top: 12rpx; color: $muted; font-size: 22rpx; }
.date { flex: 0 0 auto; padding-top: 8rpx; color: $muted; font-size: 20rpx; }
.hero { position: relative; margin-top: 40rpx; padding: 44rpx 36rpx; overflow: hidden; border-radius: 34rpx; background: linear-gradient(135deg, #64775d, #84957a); color: white; box-shadow: 0 20rpx 54rpx rgba(57, 75, 51, .12); }
.hero::after { content: ''; position: absolute; right: -70rpx; bottom: -100rpx; width: 320rpx; height: 320rpx; border-radius: 50%; border: 34rpx solid rgba(255,255,255,.08); }
.hero-kicker { display: block; color: #e9efe5; font-size: 18rpx; letter-spacing: 2rpx; }
.hero-title { display: block; margin-top: 20rpx; font-size: 40rpx; font-weight: 700; line-height: 1.45; white-space: pre-line; }
.hero-copy { display: block; max-width: 520rpx; margin-top: 16rpx; color: #edf1ea; font-size: 21rpx; line-height: 1.65; }
.hero-actions { position: relative; z-index: 1; display: flex; gap: 14rpx; margin-top: 32rpx; }
.button { min-height: 70rpx; margin: 0; padding: 0 28rpx; border-radius: 18rpx; font-size: 22rpx; line-height: 70rpx; }
.primary { background: white; color: $sage-dark; }
.light { background: rgba(255,255,255,.14); color: white; }
.metrics { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14rpx; margin-top: 24rpx; }
.metric { min-width: 0; padding: 22rpx 16rpx; display: flex; align-items: center; gap: 14rpx; border: 1rpx solid $line; border-radius: 24rpx; background: white; }
.metric-icon, .tool-icon { width: 62rpx; height: 62rpx; flex: 0 0 62rpx; display: flex; align-items: center; justify-content: center; border-radius: 20rpx; background: $sage-soft; color: $sage-dark; font-size: 26rpx; }
.peach { background: $peach-soft; }.pink { background: $pink-soft; }
.metric-label, .metric-value { display: block; }.metric-label { color: $muted; font-size: 18rpx; }.metric-value { margin-top: 4rpx; color: $ink; font-size: 32rpx; font-weight: 700; }.metric-value small { font-size: 17rpx; font-weight: 400; }
.section-head { margin-top: 42rpx; }.section-title, .section-copy { display: block; }.section-title { color: $ink; font-size: 26rpx; font-weight: 700; }.section-copy { margin-top: 6rpx; color: $muted; font-size: 20rpx; }
.tools { display: grid; gap: 14rpx; margin-top: 20rpx; }
.tool { width: 100%; min-height: 118rpx; margin: 0; padding: 20rpx 22rpx; display: flex; align-items: center; gap: 20rpx; border: 1rpx solid $line; border-radius: 26rpx; background: white; text-align: left; line-height: normal; }
.tool-title, .tool-copy { display: block; }.tool-title { color: $ink; font-size: 22rpx; font-weight: 700; }.tool-copy { margin-top: 8rpx; color: $muted; font-size: 19rpx; }.arrow { margin-left: auto; color: #a3a89e; font-size: 36rpx; }
@media (max-width: 420px) { .metrics { gap: 8rpx; }.metric { padding: 18rpx 10rpx; gap: 8rpx; }.metric-icon { width: 50rpx; height: 50rpx; flex-basis: 50rpx; }.date { display: none; } }
</style>

