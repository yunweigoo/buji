<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import FabricEntryModal from '@/components/FabricEntryModal.vue'
import { useWorkspaceStore } from '@/stores/workspace'

const workspace = useWorkspaceStore()
const showEntry = ref(false)
const keyword = ref('')
const material = ref('全部材质')
const materials = computed(() => ['全部材质', ...new Set(workspace.fabrics.map((item) => item.material).filter(Boolean))])
const filtered = computed(() => workspace.fabrics.filter((item) => {
  const haystack = [item.name, item.material, item.pattern, ...(item.styleTags ?? [])].join(' ').toLowerCase()
  return haystack.includes(keyword.value.trim().toLowerCase()) && (material.value === '全部材质' || item.material === material.value)
}))

onMounted(async () => { try { await workspace.loadWorkspace() } catch { uni.showToast({ title: '布料加载失败', icon: 'none' }) } })
function stockText(fabric: typeof workspace.fabrics[number]) {
  const pieces = fabric.pieces ?? []
  const total = pieces.reduce((sum, piece) => sum + (piece.lengthMm ?? 0), 0) / 1000
  return total > 0 ? `${pieces.length} 块 · 约 ${total.toFixed(1)} 米` : `${pieces.length} 块库存`
}
</script>

<template>
  <view class="page">
    <view class="heading"><view><text class="eyebrow">FABRIC LIBRARY</text><text class="title">我的布料收藏</text><text class="copy">按花色收藏，每一块布都有自己的故事。</text></view><button class="add" @click="showEntry = true">＋ 记录布料</button></view>
    <view class="toolbar"><view class="search"><text>⌕</text><input v-model="keyword" placeholder="搜索名称、材质、花纹或标签" /></view><scroll-view scroll-x class="filters"><view class="filter-row"><button v-for="item in materials" :key="item" :class="['chip', { active: material === item }]" @click="material = item">{{ item }}</button></view></scroll-view></view>
    <view v-if="workspace.loading" class="state">正在翻开布料册…</view><view v-else-if="workspace.error" class="state error">{{ workspace.error }}</view>
    <view v-else-if="filtered.length" class="grid">
      <view v-for="fabric in filtered" :key="fabric.id" class="card">
        <view class="cover"><image v-if="fabric.coverUrl" :src="fabric.coverUrl" mode="aspectFill" /><view v-else class="placeholder">✿</view><view class="name-band"><text>{{ fabric.name }}</text><text :class="['stock', { empty: !fabric.pieces?.length }]">{{ fabric.pieces?.length ? '有库存' : '待入库' }}</text></view></view>
        <view class="card-body"><view v-if="fabric.colors?.length" class="colors"><view v-for="color in fabric.colors" :key="`${color.name}-${color.hex}`" class="color"><view class="dot" :style="{ background: color.hex }" /><text>{{ color.name }}</text></view></view><text class="meta">{{ [fabric.pattern, fabric.material].filter(Boolean).join(' · ') || '待补充花纹与材质' }}</text><view class="rule" /><view class="stock-row"><text>{{ stockText(fabric) }}</text><text>{{ fabric.pieces?.[0]?.storageLocation || '未填写位置' }}</text></view></view>
      </view>
    </view>
    <view v-else class="empty"><text class="empty-icon">▦</text><text class="empty-title">{{ keyword ? '没有找到相符的布料' : '布料册还是空的' }}</text><text class="empty-copy">{{ keyword ? '换个关键词试试。' : '从记录手边第一款布开始吧。' }}</text></view>
    <FabricEntryModal :open="showEntry" @close="showEntry = false" />
  </view>
</template>

<style scoped lang="scss">
.page { min-height:100vh; padding:42rpx 32rpx 120rpx; background:$paper }.heading { display:flex; align-items:flex-start; justify-content:space-between; gap:24rpx }.eyebrow,.title,.copy,.meta,.empty-title,.empty-copy { display:block }.eyebrow { color:$sage; font-size:18rpx; font-weight:700; letter-spacing:2rpx }.title { margin-top:10rpx; color:$ink; font-size:42rpx; font-weight:700 }.copy { margin-top:10rpx; color:$muted; font-size:21rpx }.add { flex:0 0 auto; margin:0; padding:0 24rpx; border-radius:18rpx; background:$sage; color:white; font-size:21rpx; line-height:66rpx }
.toolbar { margin-top:30rpx }.search { height:76rpx; display:flex; align-items:center; gap:14rpx; padding:0 22rpx; border:1rpx solid $line; border-radius:22rpx; background:white; color:$muted }.search input { flex:1; font-size:21rpx }.filters { margin-top:16rpx; white-space:nowrap }.filter-row { display:flex; gap:12rpx }.chip { flex:0 0 auto; margin:0; padding:0 22rpx; border:1rpx solid $line; border-radius:999rpx; background:white; color:$muted; font-size:19rpx; line-height:56rpx }.chip.active { border-color:$sage; background:$sage-soft; color:$sage-dark }
.grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:20rpx; margin-top:28rpx }.card { overflow:hidden; border:1rpx solid $line; border-radius:28rpx; background:white; box-shadow:0 10rpx 28rpx rgba(50,58,46,.05) }.cover { position:relative; height:260rpx; background:$sage-soft }.cover image,.placeholder { width:100%; height:100% }.placeholder { display:flex; align-items:center; justify-content:center; color:$sage; font-size:68rpx }.name-band { position:absolute; left:0; right:0; bottom:0; padding:36rpx 20rpx 18rpx; display:flex; align-items:end; justify-content:space-between; gap:10rpx; color:white; font-size:25rpx; font-weight:700; background:linear-gradient(transparent,rgba(30,34,28,.72)) }.stock { flex:0 0 auto; padding:6rpx 12rpx; border-radius:999rpx; background:rgba(234,244,229,.92); color:$sage-dark; font-size:16rpx; font-weight:600 }.stock.empty { background:rgba(245,240,228,.94); color:#806f52 }.card-body { padding:18rpx 20rpx }.colors { display:flex; flex-wrap:wrap; gap:10rpx; margin-bottom:12rpx }.color { display:flex; align-items:center; gap:7rpx; color:$muted; font-size:17rpx }.dot { width:20rpx; height:20rpx; border-radius:50%; box-shadow:0 0 0 1rpx rgba(0,0,0,.1) }.meta { color:$muted; font-size:19rpx }.rule { margin:16rpx 0; height:1rpx; background:$line }.stock-row { display:flex; justify-content:space-between; gap:12rpx; color:$muted; font-size:17rpx }
.state,.empty { margin-top:80rpx; color:$muted; text-align:center }.state.error { color:#9d5c50 }.empty-icon { display:block; color:$sage; font-size:62rpx }.empty-title { margin-top:18rpx; color:$ink; font-size:26rpx; font-weight:700 }.empty-copy { margin-top:10rpx; font-size:20rpx }
@media (max-width:420px) { .heading { display:block }.add { margin-top:24rpx }.grid { grid-template-columns:1fr }.cover { height:330rpx } }
</style>
