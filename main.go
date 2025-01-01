package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sithukyaw666/iptablelb4/handler"
)

func main() {
	router := gin.Default()
	handler, err := handler.NewiptableHandler()
	if err != nil {
		logrus.Error("Can't initialize iptables instance", err)
		return
	}
	apiGroup := router.Group("/api/v1/iptables")
	{
		apiGroup.GET("/health", handler.HealthCheck)
		apiGroup.GET("/list", handler.ListRule)
		apiGroup.POST("/update", handler.UpdateRule)
		apiGroup.POST("/add", handler.AddRule)
		apiGroup.POST("/delete", handler.DeleteRule)
	}

	router.Run()
}
