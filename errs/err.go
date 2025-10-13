package errs

import "errors"

var UnauthorizedError = errors.New("未授权")

var InvalidParamError = errors.New("参数错误")

var InternalError = errors.New("内部错误")

var InvalidVerifyCodeError = errors.New("验证码无效")
