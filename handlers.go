package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (s *server) HandlePing() gin.HandlerFunc {
	return func(context *gin.Context) {
		pong := map[string]string{"message": "pong"}
		context.JSON(http.StatusOK, pong)
	}
}

func (s *server) HandleGetTokens() gin.HandlerFunc {
	return func(context *gin.Context) {
		theToken, err := GETglstsFromDB(s.db)
		if err != nil {
			context.JSON(http.StatusBadRequest, theToken)
		} else {
			context.JSON(http.StatusOK, theToken)
		}
	}
}

func (s *server) HandleGetFreshToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		theToken, err := GetFreshToken(s.db)
		fmt.Println("handle nonexpired:", err)
		if err != nil {
			context.JSON(http.StatusNoContent, nil)
		} else {
			setTokenToUsed(s.db, theToken.Steamid)
			context.JSON(http.StatusOK, theToken)
		}
	}
}

type RetiresteamIDBody struct {
	// json tag to serialize json body
	SteamID string `json:"steamid" validate:"required, len=17"`
}

func (s *server) HandleRetireToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		//setTokenToUnused(s.db, steamID)
		body := RetiresteamIDBody{}
		// using BindJson method to serialize body with struct
		if err := context.BindJSON(&body); err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fmt.Println("body!: ", body.SteamID)
		val := validator.New()

		errs := val.Var(body.SteamID, "required,len=17")

		if errs != nil {
			fmt.Println(errs)
			return
		}

		setTokenToUnused(s.db, body.SteamID)
		context.JSON(http.StatusAccepted, &body)
	}

}

//TODO: NOT IMPLEMENTED
func (s *server) HandleUsedTokens() gin.HandlerFunc {

	return func(context *gin.Context) {
		context.JSON(http.StatusNotImplemented, nil)
	}

}

func (s *server) HandleInsertDB() gin.HandlerFunc {

	return func(context *gin.Context) {

		s.INSERTglsts()
		context.JSON(http.StatusOK, "success")

	}
}
