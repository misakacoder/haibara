package ro

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"haibara/database"
	"haibara/errs"
	"haibara/model"
	"haibara/util"
	"reflect"
	"strings"
)

type LoggerRO struct {
	Level string `form:"level" binding:"required" msg:"日志级别不能为空"`
}

type IdRO struct {
	ID uint `uri:"id" binding:"required" msg:"id不能为空"`
}

type UserRO struct {
	Username string `form:"username" json:"username" binding:"required" msg:"用户名不能为空"`
	Password string `form:"password" json:"password" binding:"required" msg:"密码不能为空"`
}

type UserCreateRO struct {
	UserRO
	Nickname string `form:"nickname" json:"nickname" binding:"required" msg:"昵称不能为空"`
}

type UserUpdateRO struct {
	Nickname string `form:"nickname" json:"nickname"`
	Password string `form:"password" json:"password"`
	Enabled  bool   `form:"enabled" json:"enabled"`
}

func ValidateUri[T any](context *gin.Context, ro T) T {
	if err := context.BindUri(ro); err != nil {
		validationError(ro, err)
	}
	return ro
}

func ValidateRO[T any](context *gin.Context, ro T) T {
	if err := context.Bind(ro); err != nil {
		validationError(ro, err)
	}
	return ro
}

func ValidateUserRO(ro *UserRO) (*model.User, error) {
	user := &model.User{}
	if err := db.GORM.Where(&model.User{Username: ro.Username}).First(user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.UserNotFoundError
	}
	if util.MD5(ro.Password) != user.Password {
		return nil, errs.PasswordError
	}
	return user, nil
}

func validationError(ro any, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var msgs []string
		objectElement := reflect.TypeOf(ro).Elem()
		for _, fieldError := range validationErrors {
			if field, ok := objectElement.FieldByName(fieldError.Field()); ok {
				msg := field.Tag.Get("msg")
				if msg != "" {
					msgs = append(msgs, msg)
				}
			}
		}
		if len(msgs) > 0 {
			panic(strings.Join(msgs, "、"))
		}
	}
	panic(err)
}
