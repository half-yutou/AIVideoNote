package errcode

const (
	Success          = 0
	BadRequest       = 400
	NotFound         = 404
	InternalError    = 500

	PlatformNotSupported = 1001
	InvalidVideoURL      = 1002
	TaskNotFound         = 1003
	TaskAlreadyExists    = 1004
	DownloadFailed       = 1005
	TranscribeFailed     = 1006
	LLMCallFailed        = 1007
	ProviderNotFound     = 1008
	ProviderDisabled     = 1009
	ChatIndexNotReady    = 1010
)

var messages = map[int]string{
	PlatformNotSupported: "不支持的视频平台",
	InvalidVideoURL:      "无效的视频链接",
	TaskNotFound:         "任务不存在",
	TaskAlreadyExists:    "任务已存在",
	DownloadFailed:       "视频下载失败",
	TranscribeFailed:     "语音转录失败",
	LLMCallFailed:        "大模型调用失败",
	ProviderNotFound:     "模型供应商不存在",
	ProviderDisabled:     "模型供应商已禁用",
	ChatIndexNotReady:    "向量索引未就绪，请先构建索引",
}

func Message(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
