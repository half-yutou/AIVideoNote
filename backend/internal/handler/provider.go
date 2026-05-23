package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/errcode"
	"github.com/aivideonote/backend/internal/pkg/response"
	"github.com/aivideonote/backend/internal/service/provider"
)

type ProviderHandler struct {
	svc *provider.Service
}

func NewProviderHandler(svc *provider.Service) *ProviderHandler {
	return &ProviderHandler{svc: svc}
}

type providerRequest struct {
	Name    string `json:"name" binding:"required"`
	ApiKey  string `json:"api_key" binding:"required"`
	BaseURL string `json:"base_url" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Logo    string `json:"logo"`
}

type providerUpdateRequest struct {
	Name    *string `json:"name"`
	ApiKey  *string `json:"api_key"`
	BaseURL *string `json:"base_url"`
	Type    *string `json:"type"`
	Logo    *string `json:"logo"`
	Enabled *bool   `json:"enabled"`
}

func (h *ProviderHandler) Create(c *gin.Context) {
	var req providerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	p := &model.LLMProvider{
		Name:    req.Name,
		ApiKey:  req.ApiKey,
		BaseURL: req.BaseURL,
		Type:    req.Type,
		Logo:    req.Logo,
		Enabled: true,
	}

	if err := h.svc.Create(p); err != nil {
		response.InternalError(c, "创建提供商失败")
		return
	}

	response.Success(c, p)
}

func (h *ProviderHandler) GetAll(c *gin.Context) {
	providers, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, "获取提供商列表失败")
		return
	}

	response.Success(c, providers)
}

func (h *ProviderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ProviderNotFound, errcode.Message(errcode.ProviderNotFound))
		return
	}

	response.Success(c, p)
}

func (h *ProviderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	_, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ProviderNotFound, errcode.Message(errcode.ProviderNotFound))
		return
	}

	var req providerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ApiKey != nil {
		updates["api_key"] = *req.ApiKey
	}
	if req.BaseURL != nil {
		updates["base_url"] = *req.BaseURL
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Logo != nil {
		updates["logo"] = *req.Logo
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if len(updates) == 0 {
		response.BadRequest(c, "请至少填写一个参数")
		return
	}

	if err := h.svc.Update(id, updates); err != nil {
		response.InternalError(c, "更新提供商失败")
		return
	}

	p, _ := h.svc.GetByID(id)
	response.Success(c, p)
}

func (h *ProviderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ProviderNotFound, errcode.Message(errcode.ProviderNotFound))
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.InternalError(c, "删除提供商失败")
		return
	}

	response.SuccessMsg(c, "删除成功")
}

func (h *ProviderHandler) ListModels(c *gin.Context) {
	id := c.Query("provider_id")
	if id == "" {
		response.BadRequest(c, "请指定 provider_id")
		return
	}

	p, err := h.svc.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ProviderNotFound, errcode.Message(errcode.ProviderNotFound))
		return
	}

	if !p.Enabled {
		response.Error(c, http.StatusBadRequest, errcode.ProviderDisabled, errcode.Message(errcode.ProviderDisabled))
		return
	}

	models, err := h.svc.ListModels(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.InternalError, "获取模型列表失败: "+err.Error())
		return
	}

	response.Success(c, models)
}

func (h *ProviderHandler) TestConnection(c *gin.Context) {
	id := c.Param("id")
	if _, err := h.svc.GetByID(id); err != nil {
		response.Error(c, http.StatusNotFound, errcode.ProviderNotFound, errcode.Message(errcode.ProviderNotFound))
		return
	}

	if err := h.svc.TestConnection(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InternalError, "连接失败: "+err.Error())
		return
	}

	response.SuccessMsg(c, "连接成功")
}
