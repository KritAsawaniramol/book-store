package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthCheckRes struct {
	App    string `json:"app"`
	Status string `json:"status"`
}

func (g *ginServer) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, healthCheckRes{
		App:    g.cfg.App.Name,
		Status: "OK",
	})
}
