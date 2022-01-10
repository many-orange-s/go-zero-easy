package errorx

/*
这个里面放的都是api里面的err还有返回给前端的错误
*/
const ()

// ErrInt 自定义错误码
type ErrInt int

type CodeError struct {
	Code ErrInt `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErrorResponse 返回码
type CodeErrorResponse struct {
	Code ErrInt `json:"code"`
	Msg  string `json:"msg"`
}

// Error 只有实现这个方法才能继承Error
func (c *CodeError) Error() string {
	return c.Msg
}

// NewCodeErr 可以说是在系统里面放转动 后面会转换到Response里面
func NewCodeErr(code ErrInt, msg string) error {
	return &CodeError{code, msg}
}

func (c *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		c.Code,
		c.Msg,
	}
}
