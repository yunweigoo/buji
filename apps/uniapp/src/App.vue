<script setup lang="ts">
import { onLaunch } from '@dcloudio/uni-app'

onLaunch(() => {
  console.info('[布迹] application launched')
})
</script>

<style lang="scss">
page {
  min-height: 100%;
  background: $paper;
  color: $ink;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

view,
text,
button,
input,
textarea,
image {
  box-sizing: border-box;
}

button::after {
  border: 0;
}
</style>

