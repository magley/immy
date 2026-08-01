package handler

import (
	util "immy-api/util"
	"net/http"

	"github.com/gin-gonic/gin"

	"immy-api/service"
)

type ConfigHandler struct {
	ConfigService *service.ConfigService
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	res, err := h.ConfigService.GetConfig()

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "GET_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ConfigHandler) SetPostingEnabled(c *gin.Context) {
	enabled, ok := util.QueryIntSafe(c, "enabled", "Enabled flag")
	if !ok {
		return
	}

	res, err := h.ConfigService.SetPostingEnabled(enabled != 0)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}
