package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *server) HandlePing() gin.HandlerFunc {
	return func(context *gin.Context) {
		pong := map[string]string{"message": "pong"}
		context.JSON(http.StatusOK, pong)
	}
}

func (s *server) HandleNewToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, s.getRandomToken(s.gslList))
	}

}
