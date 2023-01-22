package routes

import (
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
	getRoutes()
	router.Run(":5001")
}

func getRoutes() {
	v1 := router.Group("/api")
	addUserRoutes(v1)
}
