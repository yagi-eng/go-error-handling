package myerror

import (
	"fmt"

	"go.uber.org/zap"
)

type MyError struct {
	Code       string
	Msg        string
	StackTrace string
}

// Error error interfaceを実装
func (me *MyError) Error() string {
	return fmt.Sprintf("my error: code[%s], message[%s]", me.Code, me.Msg)
}

// New コンストラクタ
func New(code string, msg string) *MyError {
	stack := zap.Stack("").String
	return &MyError{
		Code:       code,
		Msg:        msg,
		StackTrace: stack,
	}
}

func (me *MyError) WrapMessage(msg string) {
	me.Msg = fmt.Sprintf("%s %s", msg, me.Msg)
}
