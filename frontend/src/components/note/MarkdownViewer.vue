<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

const props = defineProps<{ markdown: string }>()

const html = computed(() => marked.parse(props.markdown) as string)
</script>

<template>
  <div class="markdown-viewer">
    <div class="markdown-body" v-html="html" />
  </div>
</template>

<style scoped>
.markdown-viewer {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}
.markdown-body {
  padding: 24px 32px;
  line-height: 1.8;
  overflow-y: auto;
}
</style>

<style>
.markdown-body h1 { font-size: 24px; margin: 16px 0 12px; }
.markdown-body h2 { font-size: 20px; margin: 16px 0 10px; border-bottom: 1px solid #eee; padding-bottom: 6px; }
.markdown-body h3 { font-size: 17px; margin: 14px 0 8px; }
.markdown-body p { margin: 8px 0; }
.markdown-body ul, .markdown-body ol { padding-left: 24px; margin: 8px 0; }
.markdown-body li { margin: 4px 0; }
.markdown-body code {
  background: #f4f4f4;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 13px;
}
.markdown-body pre {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 12px 0;
}
.markdown-body pre code { background: none; padding: 0; }
.markdown-body blockquote {
  border-left: 3px solid #1a1a2e;
  padding-left: 16px;
  color: #666;
  margin: 12px 0;
}
.markdown-body table {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}
.markdown-body th, .markdown-body td {
  border: 1px solid #ddd;
  padding: 8px 12px;
  text-align: left;
}
.markdown-body th { background: #f5f5f5; }
</style>
