package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yagi-eng/go-error-handling/myerror"
	"go.uber.org/zap"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	he, ok := err.(*echo.HTTPError)
	if !ok {
		// echo.NewHTTPErrorを使っていないとここに入る
		// e.g. panicなど
		zap.S().Errorf("Unknown error: %v", err)
		c.JSON(http.StatusInternalServerError, "panicとかでした")
		return
	}

	httpCode := he.Code
	switch err := he.Message.(type) {
	case error:
		// controllerでerror型をreturnするとここに入る
		switch {
		case httpCode >= 500:
			zap.S().Errorf("Server error: %v", err)
			if me, ok := err.(*myerror.MyError); ok {
				fmt.Print(me.StackTrace)
			}
		case httpCode >= 400:
			zap.S().Infof("Client error: %v", err)
		}
		c.JSON(httpCode, "handlingされたerrorでした")
	case string:
		// echoでエラーはが発生するとここに入る
		// e.g. 存在しないURLにアクセス
		zap.S().Errorf("Echo HTTP error: %v", he)
		c.JSON(http.StatusInternalServerError, "echoのerrorでした")
	default:
		// 通常到達しない
		zap.S().Errorf("Unknown HTTP error: %v", he)
		c.JSON(http.StatusInternalServerError, "不明なerrorです")
	}
}
