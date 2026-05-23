package llm

import "strings"

func BuildVideoSummary(transcript string, metaTitle, metaDescription string, style string) string {
	var styleGuide string
	switch style {
	case "academic":
		styleGuide = "请以学术风格撰写，使用正式、客观的语言，重点提取论点和论证过程。"
	case "casual":
		styleGuide = "请以口语化风格撰写，语言轻松自然，像朋友聊天一样。"
	case "keypoint":
		styleGuide = "请以要点提取为主，简洁列出核心要点，每点一句话。"
	default:
		styleGuide = "请以清晰的结构化风格撰写。"
	}

	var titleSection string
	if metaTitle != "" {
		titleSection = "\n视频标题：" + metaTitle
	}

	var descSection string
	if metaDescription != "" {
		descSection = "\n视频简介：" + strings.TrimSpace(metaDescription)
		if len(descSection) > 500 {
			descSection = descSection[:500]
		}
	}

	return `你是一个专业的视频内容总结助手。请根据以下视频转录文本，生成一份结构化的 Markdown 笔记。

` + styleGuide + `

要求：
1. 使用 Markdown 格式
2. 笔记内容的开头必须用 H1 标题（# 标题）的格式，为这篇笔记生成一个简洁的中文标题（15字以内），标题应概括视频核心主题
3. 使用 H2（##）标记主要章节，H3（###）标记子章节
4. 在关键位置可以插入 [screenshot:MM:SS] 标记来表示需要插入截图的时间点
5. 在关键概念处可以插入 [link:MM:SS] 标记来建立原片跳转链接
6. 提取关键数据和结论

` + titleSection + descSection + `

--- 转录文本 ---
` + transcript
}

func BuildChatContext(noteContent string, transcriptSnippets []string, question string) string {
	var sb strings.Builder

	sb.WriteString("基于以下笔记和转录内容，回答用户问题。\n\n")

	if noteContent != "" {
		sb.WriteString("## 笔记内容\n")
		sb.WriteString(noteContent)
		sb.WriteString("\n\n")
	}

	if len(transcriptSnippets) > 0 {
		sb.WriteString("## 相关原文片段\n")
		for _, chunk := range transcriptSnippets {
			sb.WriteString("- ")
			sb.WriteString(chunk)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 用户问题\n")
	sb.WriteString(question)

	return sb.String()
}
