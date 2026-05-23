<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="layout">
    <header class="header">
      <div class="header-left">
        <router-link to="/" class="logo">VideoNote</router-link>
        <span class="subtitle">AI 视频笔记</span>
      </div>
      <nav class="header-nav">
        <router-link to="/" class="nav-link">首页</router-link>
        <router-link to="/settings" class="nav-link">设置</router-link>
      </nav>
    </header>
    <main class="main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
.header {
  height: 52px;
  background: #1a1a2e;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
}
.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}
.logo {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  text-decoration: none;
}
.logo:hover {
  color: #64b5f6;
}
.subtitle {
  font-size: 13px;
  color: #aaa;
}
.header-nav {
  display: flex;
  gap: 16px;
}
.nav-link {
  color: #ccc;
  text-decoration: none;
  font-size: 14px;
  padding: 4px 12px;
  border-radius: 4px;
}
.nav-link:hover,
.nav-link.active {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}
.main {
  flex: 1;
  overflow: hidden;
}
</style>
