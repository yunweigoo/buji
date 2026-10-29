<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useWorkspaceStore } from '@/stores/workspace'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; saved: [] }>()
const workspace = useWorkspaceStore()
const saving = ref(false)
const imagePath = ref('')
const form = reactive({
  name: '', material: '棉', pattern: '', tags: '', colorName: '', colorHex: '#73836d',
  lengthM: '', widthCm: '', location: '', channel: '', priceYuan: '', purchaseDate: '', memo: '',
})
const canSave = computed(() => Boolean(form.name.trim()) && !saving.value)

watch(() => props.open, (open) => {
  if (!open) return
  imagePath.value = ''
  Object.assign(form, { name: '', material: '棉', pattern: '', tags: '', colorName: '', colorHex: '#73836d', lengthM: '', widthCm: '', location: '', channel: '', priceYuan: '', purchaseDate: '', memo: '' })
})

function chooseImage() {
  uni.chooseImage({ count: 1, sizeType: ['compressed'], sourceType: ['album', 'camera'], success: ({ tempFilePaths }) => { imagePath.value = tempFilePaths[0] ?? '' } })
}

function optionalInt(value: string, multiplier: number) {
  const parsed = Number(value)
  return Number.isFinite(parsed) && parsed > 0 ? Math.round(parsed * multiplier) : undefined
}

async function save() {
  if (!canSave.value) return
  saving.value = true
  try {
    await workspace.createFabric({
      name: form.name.trim(), material: form.material.trim(), pattern: form.pattern.trim(),
      styleTags: form.tags.split(/[，,]/).map((item) => item.trim()).filter(Boolean),
      colors: form.colorName.trim() ? [{ name: form.colorName.trim(), hex: form.colorHex }] : [],
      memo: form.memo.trim(),
      piece: {
        label: '首块库存', lengthMm: optionalInt(form.lengthM, 1000), widthMm: optionalInt(form.widthCm, 10),
        storageLocation: form.location.trim(), purchaseChannel: form.channel.trim(), purchaseDate: form.purchaseDate || undefined,
        totalPriceCent: optionalInt(form.priceYuan, 100), memo: '',
      },
    }, imagePath.value || undefined)
    uni.showToast({ title: '已存入布料库', icon: 'success' })
    emit('saved')
    emit('close')
  } catch (error) {
    uni.showToast({ title: error instanceof Error ? error.message : '保存失败', icon: 'none' })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <view v-if="open" class="overlay" @click.self="emit('close')">
    <view class="sheet">
      <view class="sheet-head"><view><text class="kicker">NEW FABRIC</text><text class="sheet-title">记录一款布料</text></view><button class="close" @click="emit('close')">×</button></view>
      <scroll-view scroll-y class="body">
        <button class="photo" @click="chooseImage">
          <image v-if="imagePath" :src="imagePath" mode="aspectFill" class="preview" />
          <view v-else class="photo-empty"><text class="camera">📷</text><text class="photo-title">拍照或选择布料照片</text><text class="photo-copy">照片会作为布料卡片封面</text></view>
          <text v-if="imagePath" class="replace">更换照片</text>
        </button>
        <view class="section"><text class="section-title">识别信息</text><text class="section-copy">先确认基础特征，后续会补充自动识别。</text></view>
        <view class="field"><text class="label">布料名称 *</text><input v-model="form.name" placeholder="例如：蓝色碎花棉布" /></view>
        <view class="row"><view class="field"><text class="label">材质</text><input v-model="form.material" placeholder="棉 / 麻 / 毛" /></view><view class="field"><text class="label">花纹</text><input v-model="form.pattern" placeholder="碎花 / 格纹" /></view></view>
        <view class="field"><text class="label">风格标签</text><input v-model="form.tags" placeholder="日系，复古（逗号分隔）" /></view>
        <view class="row"><view class="field"><text class="label">主色名称</text><input v-model="form.colorName" placeholder="鼠尾草绿" /></view><view class="field color-field"><text class="label">色值</text><input v-model="form.colorHex" /><view class="color-dot" :style="{ background: form.colorHex }" /></view></view>
        <view class="section"><text class="section-title">实物布块</text><text class="section-copy">尺寸、存放和购买信息会进入库存记录。</text></view>
        <view class="row"><view class="field"><text class="label">长度（米）</text><input v-model="form.lengthM" type="digit" placeholder="2.5" /></view><view class="field"><text class="label">幅宽（厘米）</text><input v-model="form.widthCm" type="digit" placeholder="145" /></view></view>
        <view class="row"><view class="field"><text class="label">存放位置</text><input v-model="form.location" placeholder="布柜 A2" /></view><view class="field"><text class="label">购买渠道</text><input v-model="form.channel" placeholder="门店 / 淘宝" /></view></view>
        <view class="row"><view class="field"><text class="label">总价（元）</text><input v-model="form.priceYuan" type="digit" placeholder="68" /></view><view class="field"><text class="label">购买日期</text><input v-model="form.purchaseDate" type="text" placeholder="2026-10-29" /></view></view>
        <view class="field"><text class="label">备注</text><textarea v-model="form.memo" placeholder="适合做什么、手感、来源故事……" /></view>
      </scroll-view>
      <view class="footer"><button class="cancel" @click="emit('close')">稍后再说</button><button class="save" :disabled="!canSave" @click="save">{{ saving ? '正在保存…' : '保存到布料库' }}</button></view>
    </view>
  </view>
</template>

<style scoped lang="scss">
.overlay { position: fixed; z-index: 99; inset: 0; display: flex; justify-content: center; align-items: flex-end; background: rgba(45, 48, 41, .46); backdrop-filter: blur(6rpx); }
.sheet { width: min(920rpx, 100%); height: 92vh; display: flex; flex-direction: column; overflow: hidden; border-radius: 36rpx 36rpx 0 0; background: $paper; box-shadow: 0 -24rpx 80rpx rgba(34, 39, 31, .18); }
.sheet-head { padding: 30rpx 32rpx 22rpx; display: flex; justify-content: space-between; align-items: center; border-bottom: 1rpx solid $line; background: rgba(255,255,255,.76); }.kicker,.sheet-title { display:block }.kicker { color:$sage; font-size:17rpx; font-weight:700; letter-spacing:2rpx }.sheet-title { margin-top:6rpx; color:$ink; font-size:30rpx; font-weight:700 }.close { width:58rpx; height:58rpx; margin:0; padding:0; border-radius:50%; background:#eef0eb; color:$muted; font-size:38rpx; line-height:54rpx }
.body { flex: 1; min-height: 0; padding: 26rpx 32rpx; box-sizing: border-box; }.photo { position:relative; width:100%; height:300rpx; margin:0; padding:0; overflow:hidden; border:2rpx dashed #b9c2b4; border-radius:28rpx; background:$sage-soft; line-height:normal }.preview { width:100%; height:100% }.photo-empty { height:100%; display:flex; flex-direction:column; align-items:center; justify-content:center }.camera { font-size:44rpx }.photo-title { margin-top:14rpx; color:$ink; font-size:23rpx; font-weight:700 }.photo-copy { margin-top:8rpx; color:$muted; font-size:18rpx }.replace { position:absolute; right:18rpx; bottom:18rpx; padding:10rpx 16rpx; border-radius:16rpx; background:rgba(255,255,255,.9); color:$ink; font-size:18rpx }
.section { margin:34rpx 0 18rpx }.section-title,.section-copy,.label { display:block }.section-title { color:$ink; font-size:24rpx; font-weight:700 }.section-copy { margin-top:6rpx; color:$muted; font-size:18rpx }.row { display:grid; grid-template-columns:1fr 1fr; gap:18rpx }.field { position:relative; margin-bottom:18rpx }.label { margin-bottom:10rpx; color:$muted; font-size:18rpx }.field input,.field textarea { width:100%; box-sizing:border-box; border:1rpx solid $line; border-radius:18rpx; background:white; color:$ink; font-size:21rpx }.field input { height:72rpx; padding:0 20rpx }.field textarea { height:150rpx; padding:18rpx 20rpx }.color-field input { padding-right:62rpx }.color-dot { position:absolute; right:18rpx; bottom:20rpx; width:30rpx; height:30rpx; border:4rpx solid white; border-radius:50%; box-shadow:0 0 0 1rpx $line }
.footer { padding:20rpx 32rpx calc(20rpx + env(safe-area-inset-bottom)); display:flex; gap:14rpx; border-top:1rpx solid $line; background:white }.footer button { margin:0; border-radius:18rpx; font-size:21rpx; line-height:74rpx }.cancel { flex:0 0 180rpx; background:#f1f1ed; color:$muted }.save { flex:1; background:$sage; color:white }.save[disabled] { opacity:.45 }
@media (min-width: 700px) { .overlay { align-items:center; padding:40rpx }.sheet { height:min(920rpx, 90vh); border-radius:36rpx } }
</style>
