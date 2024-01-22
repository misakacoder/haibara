package model

import (
	"fmt"
	db "haibara/database"
	"haibara/definition"
	"haibara/util"
)

const tableOptionsFormat = "engine=InnoDB default charset=utf8mb4 collate=utf8mb4_bin comment='%s'"

type User struct {
	ID        uint                 `gorm:"column:id;type:int;autoIncrement;primaryKey;comment:主键" json:"-"`
	Username  string               `gorm:"column:username;type:varchar(64);not null;uniqueIndex:udx_username;comment:用户名" json:"username"`
	Nickname  string               `gorm:"column:nickname;type:varchar(64);not null;comment:昵称" json:"nickname"`
	Password  string               `gorm:"column:password;type:varchar(255);not null;comment:密码" json:"-"`
	Enabled   bool                 `gorm:"column:enabled;type:tinyint(1);not null;default:1;comment:是否启用" json:"enable"`
	CreatedAt *definition.DateTime `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"`
	UpdatedAt *definition.DateTime `gorm:"column:updated_at;type:datetime;null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
}

func (user *User) Roles() []definition.Role {
	var roles []definition.Role
	db.GORM.Model(&UserRole{}).Where(&UserRole{UserID: user.ID}).Pluck("role", &roles)
	return roles
}

type UserRole struct {
	ID     uint            `gorm:"column:id;type:int;autoIncrement;primaryKey;comment:主键" json:"-"`
	UserID uint            `gorm:"column:user_id;type:int;not null;uniqueIndex:udx_user_id_role;comment:用户id"`
	Role   definition.Role `gorm:"column:role;type:int;not null;uniqueIndex:udx_user_id_role;comment:角色"`
}

func FirstOrCreate() {
	tables := map[string]any{
		"用户表":   &User{},
		"用户角色表": &UserRole{},
	}
	for k, v := range tables {
		_ = db.GORM.Set("gorm:table_options", fmt.Sprintf(tableOptionsFormat, k)).AutoMigrate(v)
	}

	users := map[*User][]*UserRole{
		&User{Username: "haibara", Nickname: "灰原哀", Password: util.MD5("haibara"), Enabled: true}: {
			&UserRole{Role: definition.ADMIN},
			&UserRole{Role: definition.USER},
		},
		&User{Username: "gin", Nickname: "琴酒", Password: util.MD5("gin"), Enabled: true}: {
			&UserRole{Role: definition.USER},
		},
	}
	for user, roles := range users {
		db.GORM.Where(&User{Username: user.Username}).FirstOrCreate(user)
		id := user.ID
		for _, role := range roles {
			role.UserID = id
			db.GORM.Where(role).FirstOrCreate(role)
		}
	}

	init := []any{}
	for _, data := range init {
		db.GORM.Where(data).FirstOrCreate(data)
	}
}
