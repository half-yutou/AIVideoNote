package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/group"
)

type GroupHandler struct {
	svc *group.Service
}

func NewGroupHandler(svc *group.Service) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) List(c *gin.Context) {
	groups, err := h.svc.List()
	if err != nil {
		response.InternalError(c, "获取分组列表失败")
		return
	}
	response.Success(c, groups)
}

type createGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *GroupHandler) Create(c *gin.Context) {
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	g, err := h.svc.Create(req.Name)
	if err != nil {
		response.InternalError(c, "创建分组失败")
		return
	}
	response.Success(c, g)
}

func (h *GroupHandler) Rename(c *gin.Context) {
	id := c.Param("id")
	var req renameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Rename(id, req.Name); err != nil {
		response.InternalError(c, "重命名分组失败")
		return
	}
	response.SuccessMsg(c, "重命名成功")
}

func (h *GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "default" {
		response.Error(c, http.StatusBadRequest, 400, "不能删除默认分组")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.InternalError(c, "删除分组失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}
