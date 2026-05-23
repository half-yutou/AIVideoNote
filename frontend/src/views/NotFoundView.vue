<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <div class="not-found">
    <h1>404</h1>
    <p>页面不存在</p>
    <button @click="router.push('/')">返回首页</button>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: calc(100vh - 52px);
  color: #666;
}
.not-found h1 { font-size: 72px; margin-bottom: 8px; }
.not-found p { font-size: 18px; margin-bottom: 24px; }
.not-found button {
  padding: 10px 24px;
  background: #1a1a2e;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
</style>
