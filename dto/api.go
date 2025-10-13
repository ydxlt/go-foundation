package dto

const (
	// CodeSuccess 通用错误码
	CodeSuccess           = 0
	CodeInternalServerErr = 10000
	CodeParamInvalid      = 10001
	CodeUnauthorized      = 10002
	CodeForbidden         = 10003
	CodeNotFound          = 10004
	CodeWithMessage       = 10005
	CodeAuthInvalid       = 10006
	CodeCommonError       = 10007
	CodeBadRequest        = 10007
	CodeBadStatus         = 10008

	// CodeUserNotFound UserNotFound 用户相关
	CodeUserNotFound       = 10100
	CodeUserAlreadyExists  = 10101
	CodeInvalidPassword    = 10102
	CodeEmailCodeInvalid   = 10103
	CodeTokenExpired       = 10104
	CodeRefreshTokenExpire = 10105
	CodeUseBad             = 10106

	// CodeDataCreateFailed DataCreateFailed 数据相关
	CodeDataCreateFailed = 10200
	CodeDataUpdateFailed = 10201
	CodeDataDeleteFailed = 10202
	CodeDataQueryFailed  = 10203

	// CodePermissionDenied PermissionDenied 权限相关
	CodePermissionDenied = 10300
)

// Response 通用返回结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

// Builder builder
type Builder struct {
	resp *Response
}

// New 创建 builder
func New() *Builder {
	return &Builder{
		resp: &Response{},
	}
}

// Code 设置 Code
func (b *Builder) Code(code int) *Builder {
	b.resp.Code = code
	return b
}

// Msg 设置消息
func (b *Builder) Msg(msg string) *Builder {
	b.resp.Msg = msg
	return b
}

// Data 设置数据
func (b *Builder) Data(data interface{}) *Builder {
	b.resp.Data = data
	return b
}

// Build 构建 Response
func (b *Builder) Build() *Response {
	return b.resp
}

//
// ========== 快捷方法 ==========
//

// Success 成功（带数据）
func Success(data interface{}) *Response {
	return &Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	}
}

// Ok 成功（无数据）
func Ok() *Response {
	return &Response{
		Code: CodeSuccess,
		Msg:  "success",
	}
}

// Error 错误（自定义 code + msg）
func Error(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}

// BadRequest 常见错误
func BadRequest(msg string) *Response {
	return Error(CodeBadRequest, msg)
}

// UnauthorizedError 未授权
func UnauthorizedError() *Response {
	return Error(CodeUnauthorized, "未授权，请登录")
}

// InternalError 禁止访问
func InternalError() *Response {
	return Error(CodeInternalServerErr, "内部错误，请联系管理员")
}

// ForbiddenError 禁止访问
func ForbiddenError(msg string) *Response {
	return Error(CodeForbidden, msg)
}

// MessageError 带消息的错误（客户端Toast提示）
func MessageError(msg string) *Response {
	return Error(CodeWithMessage, msg)
}

// ParamError 带消息的错误（客户端Toast提示）
func ParamError() *Response {
	return Error(CodeParamInvalid, "参数错误")
}

// UserNoFoundError 用户不存在
func UserNoFoundError() *Response {

	return Error(CodeUserNotFound, "用户不存在")
}
func CommonError(msg string) *Response {
	return Error(CodeCommonError, msg)
}
