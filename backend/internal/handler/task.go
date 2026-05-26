package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/task"
)

type TaskHandler struct {
	svc      *task.Service
	executor *task.Executor
}

func NewTaskHandler(svc *task.Service, executor *task.Executor) *TaskHandler {
	return &TaskHandler{svc: svc, executor: executor}
}

type generateRequest struct {
	VideoURL   string   `json:"video_url" binding:"required"`
	Platform   string   `json:"platform" binding:"required"`
	Quality    string   `json:"quality"`
	ModelName  string   `json:"model_name" binding:"required"`
	ProviderID string   `json:"provider_id" binding:"required"`
	GroupID    string   `json:"group_id"`
	Format     []string `json:"format"`
	Style      string   `json:"style"`
	Extras     string   `json:"extras"`
	Link       bool     `json:"link"`
}

func (h *TaskHandler) Generate(c *gin.Context) {
	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Quality == "" {
		req.Quality = "medium"
	}
	if req.GroupID == "" {
		req.GroupID = "default"
	}

	t := &model.Task{
		ID:         uuid.New().String(),
		VideoURL:   req.VideoURL,
		Platform:   req.Platform,
		Quality:    req.Quality,
		ModelName:  req.ModelName,
		ProviderID: req.ProviderID,
		GroupID:    req.GroupID,
		Format:     marshalStrSlice(req.Format),
		Style:      req.Style,
		Extras:     req.Extras,
		Link:       req.Link,
	}

	if err := h.svc.Create(t); err != nil {
		response.InternalError(c, "创建任务失败")
		return
	}

	h.executor.Submit(t.ID)

	response.Success(c, gin.H{"task_id": t.ID})
}

func (h *TaskHandler) GetStatus(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.FindByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 1003, "任务不存在")
		return
	}

	note, _ := h.svc.FindNoteByTaskID(id)

	data := gin.H{
		"task_id":       t.ID,
		"status":        t.Status,
		"video_url":     t.VideoURL,
		"platform":      t.Platform,
		"error_message": t.ErrorMessage,
		"created_at":    t.CreatedAt,
	}

	if note != nil {
		data["markdown"] = note.Markdown
		data["note_id"] = note.ID
	}

	response.Success(c, data)
}

func (h *TaskHandler) GetList(c *gin.Context) {
	tasks, err := h.svc.FindAll()
	if err != nil {
		response.InternalError(c, "获取任务列表失败")
		return
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	groupMap := make(map[string]string)
	groups, _ := h.svc.FindAllGroups()
	for _, g := range groups {
		groupMap[g.ID] = g.Name
	}

	type taskItem struct {
		ID           string           `json:"id"`
		Status       model.TaskStatus `json:"status"`
		VideoURL     string           `json:"video_url"`
		Platform     string           `json:"platform"`
		VideoID      string           `json:"video_id"`
		Name         string           `json:"name"`
		GroupName    string           `json:"group_name"`
		ErrorMessage string           `json:"error_message"`
		CreatedAt    string           `json:"created_at"`
	}

	var items []taskItem
	for _, t := range tasks {
		groupName := groupMap[t.GroupID]
		if groupName == "" {
			groupName = "默认"
		}
		items = append(items, taskItem{
			ID:           t.ID,
			Status:       t.Status,
			VideoURL:     t.VideoURL,
			Platform:     t.Platform,
			VideoID:      t.VideoID,
			Name:         t.Name,
			GroupName:    groupName,
			ErrorMessage: t.ErrorMessage,
			CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	response.Success(c, items)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.svc.FindByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 1003, "任务不存在")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.InternalError(c, "删除失败")
		return
	}

	response.SuccessMsg(c, "删除成功")
}

type renameRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *TaskHandler) RenameTask(c *gin.Context) {
	id := c.Param("id")
	var req renameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.UpdateName(id, req.Name); err != nil {
		response.InternalError(c, "重命名失败")
		return
	}
	response.SuccessMsg(c, "重命名成功")
}

func marshalStrSlice(s []string) string {
	if len(s) == 0 {
		return "[]"
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "[]"
	}
	return string(b)
}
