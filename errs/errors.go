package errs

import "errors"

var (
	UserNotFoundError = errors.New("用户不存在")
	PasswordError     = errors.New("密码错误")
	ParseTokenError   = errors.New("解析token失败")
)
