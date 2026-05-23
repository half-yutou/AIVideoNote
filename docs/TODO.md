# AIVideoNote TODO

还未实现/未完成的功能记录。

## 高优先级

- **本地文件上传**：`POST /api/v1/upload` 后端已实现，前端 `api/upload.ts` 已封装，但 NoteForm 中 `local` 平台选项未对接上传逻辑。`UPLOAD_DIR` 环境变量已预留。
- **截图插入**：后端 pipeline 中 `generateThumbnails` 为占位实现。需要对接视频帧提取 → 截图生成 → Markdown 中 `[screenshot:MM:SS]` 标记替换。前端 NoteForm 中 `screenshot` 复选框和画质选择已移除，实现后需恢复 UI。
- **原片跳转链接**：prompt 中要求 AI 生成 `[link:MM:SS]` 标记。后处理需要将标记转换为可点击链接。前端 NoteForm 中 `link` 复选框已移除，实现后需恢复 UI。

## 中优先级

- **聊天索引状态 API**：`GET /api/v1/chat/status` 后端已路由，handler 中有占位实现，前端已移除问答模块。
- **向量嵌入/视觉理解**：`transfer/app/vision/` 和 `transfer/app/embedding/` 仅有空壳路由，预留了 `/api/v1/vision/*` 和 `/api/v1/embedding/*` 端点。

## 低优先级

- **分组删除**：前端 TaskHistory 未暴露删除分组按钮，groupStore 有 `removeGroup` 方法可用。
- **笔记导出格式**：目前仅支持 Markdown 导出，未实现 PDF/DOCX。
- **移动端适配**：三栏布局在窄屏幕上体验不佳。
