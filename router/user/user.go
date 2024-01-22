package user

import (
	"github.com/gin-gonic/gin"
	"haibara/consts"
	"haibara/database"
	"haibara/definition"
	"haibara/middleware"
	"haibara/model"
	"haibara/response"
	"haibara/ro"
	"haibara/util"
)

func LoadRouter(router *gin.Engine) {
	user := router.Group("/api/v1")
	user.GET("/login", login)

	user.Use(middleware.Auth)
	user.GET("/users", middleware.HasRole(definition.ADMIN), pageUser)
	user.POST("/user", addUser)
	user.PUT("/user/:id", updateUser)
	user.DELETE("/user/:id", deleteUser)
}

func login(context *gin.Context) {
	userRO := ro.ValidateRO(context, &ro.UserRO{})
	_, err := ro.ValidateUserRO(userRO)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	text, err := util.ToJSONString(userRO)
	if err != nil {
		response.ErrorWithMessageJSON(context, err.Error())
		return
	}
	setCookie(context, consts.LoginCookieKey, util.EncAES(text))
	response.OkWithDataJSON(context, nil)
}

func pageUser(context *gin.Context) {
	page := model.Page{}
	_ = context.ShouldBind(&page)
	pageResult := model.Paginate(page, &model.User{})
	response.OkWithDataJSON(context, pageResult)
}

func addUser(context *gin.Context) {
	userRO := ro.ValidateRO(context, &ro.UserCreateRO{})
	tx := db.GORM.Create(&model.User{Username: userRO.Username, Password: util.MD5(userRO.Password), Nickname: userRO.Nickname})
	if err := tx.Error; err != nil {
		response.ErrorWithDataJSON(context, err.Error())
	} else {
		response.Ok(context)
	}
}

func updateUser(context *gin.Context) {
	idRO := ro.ValidateUri(context, &ro.IdRO{})
	userRO := ro.ValidateRO(context, &ro.UserUpdateRO{})
	user := model.User{ID: idRO.ID, Nickname: userRO.Nickname, Enabled: userRO.Enabled}
	if userRO.Password != "" {
		user.Password = util.MD5(userRO.Password)
	}
	tx := db.GORM.Updates(user)
	if err := tx.Error; err != nil {
		response.ErrorWithDataJSON(context, err.Error())
	} else {
		response.Ok(context)
	}
}

func deleteUser(context *gin.Context) {
	id := context.Param("id")
	if err := db.GORM.Delete(&model.User{}, id).Error; err != nil {
		response.ErrorWithDataJSON(context, err.Error())
	} else {
		response.Ok(context)
	}
}

func setCookie(context *gin.Context, key string, value string) {
	ip := context.ClientIP()
	ip = util.ConditionalExpression(ip == "::1", consts.Localhost, ip)
	context.SetCookie(key, value, 30*24*60*60, "/", ip, false, true)
}
