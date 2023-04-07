package errors

type MyError struct {
	Code int
	Msg  string
	Data interface{}
}

var (
	LOGIN_UNKNOWN = NewError(20200, "用户不存在")
	LOGIN_ERROR   = NewError(20300, "账号或密码错误")
	ERROR         = NewError(40000, "操作失败")
	UNAUTHORIZED  = NewError(40100, "未登录")
	ILLEGAL_DATA  = NewError(40300, "非法参数")
	NOT_FOUND     = NewError(40400, "资源不存在")
	VALID_ERROR   = NewError(40500, "参数错误")
	INNER_ERROR   = NewError(50000, "系统发生异常")
)

func (e *MyError) Error() string {
	return e.Msg
}
func CreateError(code int, msg string, data interface{}) *MyError {
	return &MyError{
		Msg:  msg,
		Code: code,
		Data: data,
	}
}

func NewError(code int, msg string) *MyError {
	return &MyError{
		Msg:  msg,
		Code: code,
	}
}

func GetError(e *MyError, data interface{}) *MyError {
	return &MyError{
		Msg:  e.Msg,
		Code: e.Code,
		Data: data,
	}
}
