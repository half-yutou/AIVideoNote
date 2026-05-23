package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/cookie"
)

type CookieHandler struct {
	svc *cookie.Service
}

func NewCookieHandler(svc *cookie.Service) *CookieHandler {
	return &CookieHandler{svc: svc}
}

type saveCookieRequest struct {
	Platform string `json:"platform" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (h *CookieHandler) GetAll(c *gin.Context) {
	cookies, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, "获取 Cookie 列表失败")
		return
	}
	response.Success(c, cookies)
}

func (h *CookieHandler) Save(c *gin.Context) {
	var req saveCookieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	cookie, err := h.svc.Save(req.Platform, req.Content)
	if err != nil {
		response.InternalError(c, "保存 Cookie 失败")
		return
	}
	response.Success(c, cookie)
}

func (h *CookieHandler) Delete(c *gin.Context) {
	platform := c.Param("platform")
	if err := h.svc.Delete(platform); err != nil {
		response.InternalError(c, "删除 Cookie 失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}
