package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sithukyaw666/iptablelb4/handler"
)

func main() {
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.Default())

	handler, err := handler.NewiptableHandler()
	if err != nil {
		logrus.Error("Can't initialize iptables instance", err)
		return
	}
	apiGroup := router.Group("/api/v1/iptables")
	{
		apiGroup.GET("/health", handler.HealthCheck)
		apiGroup.GET("/list", handler.ListFarm)
		apiGroup.GET("/list/:farm", handler.ListFarmByName)
		apiGroup.POST("/update", handler.UpdateRule)
		apiGroup.POST("/add", handler.AddRule)
		apiGroup.POST("/delete/:farm", handler.DeleteRule)
	}

	router.Run()
}
