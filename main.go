package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yagi-eng/go-error-handling/myerror"
	"go.uber.org/zap"
)

func main() {
	// zap setting
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableStacktrace = true
	logger, _ := zapConfig.Build()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	e := echo.New()
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	e.GET("/hoge", controllerFunc())
	e.Logger.Fatal(e.Start(":8080"))
}

// controllerFunc controller的な関数
func controllerFunc() echo.HandlerFunc {
	return func(c echo.Context) error {
		str, err := serviceFunc()
		if err != nil {
			httpCode := http.StatusInternalServerError
			if me, ok := err.(*myerror.MyError); ok && me.Code == "001-001" {
				httpCode = http.StatusNotFound
			}
			return echo.NewHTTPError(httpCode, err)
		}
		return c.JSON(http.StatusOK, str)
	}
}

// serviceFunc service, usecase的な関数
// あくまで例なのでentityFunc呼ぶだけ
func serviceFunc() (string, error) {
	err := entityFunc()
	if err != nil {
		return "", err
	}
	return "success!", nil
}

// entityFunc 異なる種類のエラーを返すような関数
func entityFunc() error {
	// 例えばos.Openを呼ぶ
	file := "xxx/xxx"
	_, err := os.Open(file)
	if err != nil {
		return myerror.New("001-001", err.Error())
	}
	return myerror.New("001-999", "成功してるけど、例としてとにかくエラーを返す")
}
