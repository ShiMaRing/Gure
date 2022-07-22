package module

import (
	"strings"
)

type CustomError interface {
	Type() ErrorType //指出错误出现的模块
	Error() string   //实现Error接口
}

type SpiderError struct {
	errType   ErrorType //指出错误出现的模块
	errMsg    string    //提示错误信息
	cplErrMsg string    //完整错误信息
}

// NewSpiderError 构造方法
func NewSpiderError(errType ErrorType, errMsg string) *SpiderError {
	return &SpiderError{errType: errType,
		errMsg: strings.TrimSpace(errMsg),
	}
}

// Type 返回错误类型
func (s SpiderError) Type() ErrorType {
	return s.errType
}

//Error 返回格式化后的错误信息
func (s SpiderError) Error() string {
	if s.cplErrMsg == "" {
		s.getCplErrMsg()
	}
	return s.cplErrMsg
}

//getCplErrMsg 应当采用builder形式避免字符串拼接带来的性能影响
func (s SpiderError) getCplErrMsg() {
	builder := strings.Builder{}
	builder.WriteString("Type:")
	if s.errType == "" {
		builder.WriteString("Unknown")
	} else {
		builder.WriteString(string(s.errType))
	}

	builder.WriteString(" Msg:")
	builder.WriteString(s.errMsg)

	s.cplErrMsg = builder.String()
}
