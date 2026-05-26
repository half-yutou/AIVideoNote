<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

const props = defineProps<{ markdown: string }>()

const html = computed(() => marked.parse(props.markdown) as string)
</script>

<template>
  <div class="markdown-viewer animate-fade-in-up">
    <div class="markdown-body" v-html="html" />
  </div>
</template>

<style scoped>
.markdown-viewer {
  background: var(--color-bg-elevated);
  border-radius: var(--radius-lg);
  overflow-y: auto;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  min-height: 0;
  flex: 1;
}
.markdown-body {
  padding: var(--space-8) var(--space-8);
  line-height: var(--leading-relaxed);
  font-size: var(--text-base);
  color: var(--color-text);
}
</style>

<style>
.markdown-body h1 {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 1.5rem 0 1rem;
  color: var(--color-text);
  letter-spacing: -0.02em;
}
.markdown-body h2 {
  font-size: 1.375rem;
  font-weight: 600;
  margin: 1.5rem 0 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--color-border-light);
  color: var(--color-text);
}
.markdown-body h3 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 1.25rem 0 0.5rem;
  color: var(--color-text);
}
.markdown-body h4 {
  font-size: 1rem;
  font-weight: 600;
  margin: 1rem 0 0.5rem;
  color: var(--color-text-secondary);
}
.markdown-body p {
  margin: 0.625rem 0;
  line-height: 1.75;
}
.markdown-body ul, .markdown-body ol {
  padding-left: 1.5rem;
  margin: 0.625rem 0;
}
.markdown-body li {
  margin: 0.375rem 0;
  line-height: 1.7;
}
.markdown-body li::marker {
  color: var(--color-primary);
}
.markdown-body code {
  background: var(--color-bg-muted);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.85em;
  font-family: var(--font-mono);
  color: var(--color-primary-hover);
}
.markdown-body pre {
  background: #1e293b;
  color: #e2e8f0;
  padding: 1.25rem;
  border-radius: var(--radius-md);
  overflow-x: auto;
  margin: 1rem 0;
  font-size: 0.8125rem;
  line-height: 1.6;
}
.markdown-body pre code {
  background: none;
  padding: 0;
  color: inherit;
  font-size: inherit;
}
.markdown-body blockquote {
  border-left: 3px solid var(--color-primary);
  padding: 0.5rem 1rem;
  margin: 1rem 0;
  background: var(--color-primary-light);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  color: var(--color-text-secondary);
}
.markdown-body blockquote p {
  margin: 0.25rem 0;
}
.markdown-body table {
  border-collapse: collapse;
  width: 100%;
  margin: 1rem 0;
  font-size: 0.875rem;
}
.markdown-body th, .markdown-body td {
  border: 1px solid var(--color-border);
  padding: 0.625rem 1rem;
  text-align: left;
}
.markdown-body th {
  background: var(--color-bg-muted);
  font-weight: 600;
  color: var(--color-text-secondary);
}
.markdown-body tr:hover td {
  background: var(--color-bg-muted);
}
.markdown-body hr {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 1.5rem 0;
}
.markdown-body strong {
  font-weight: 600;
  color: var(--color-text);
}
.markdown-body a {
  color: var(--color-primary);
  text-decoration: underline;
  text-underline-offset: 2px;
}
.markdown-body a:hover {
  color: var(--color-primary-hover);
}
.markdown-body img {
  max-width: 100%;
  border-radius: var(--radius-md);
  margin: 0.5rem 0;
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
