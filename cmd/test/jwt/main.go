package main

import (
	"onx-outgoing-go/internal/handler/auth"
	"onx-outgoing-go/internal/pkg/jwt"
	"onx-outgoing-go/internal/pkg/logger"
	"onx-outgoing-go/internal/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Setup()
	jwtOpts := jwt.DefaultOptions("1234")
	jwtOpts.TokenExpiredTime = 60 * time.Second
	jwtAuth := jwt.New(jwtOpts)
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RequestInit())
	r.Use(middleware.ResponseInit())

	handler := auth.NewHandler(jwtAuth)
	handler.NewRoutes(r.Group("/api"), jwtAuth)

	err := r.Run(":8001")
	if err != nil {
		return
	}
}
