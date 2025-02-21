package middleware

import (
	"errors"
	"net/http"
	_type "onx-outgoing-go/internal/common/type"
	"onx-outgoing-go/internal/pkg/helper"
	"onx-outgoing-go/internal/pkg/jwt"
	"onx-outgoing-go/internal/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(auth jwt.IJWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		send := c.MustGet("send").(func(r *_type.Response))
		token := c.GetHeader("Authorization")
		if token == "" {
			send(helper.ParseResponse(&_type.Response{Code: http.StatusBadRequest, Message: "token not found"}))
			return
		}

		logger.Debug.Println("token", token)
		parts := strings.Split(token, " ")
		if len(parts) < 2 {
			send(helper.ParseResponse(&_type.Response{Code: http.StatusBadRequest, Message: "invalid token format", Error: errors.New("invalid token format")}))
			return
		}
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			send(helper.ParseResponse(&_type.Response{Code: http.StatusBadRequest, Message: "invalid token", Error: err}))
			return
		}

		c.Set("auth", claims)
		c.Next()
	}
}
