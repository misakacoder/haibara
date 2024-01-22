package definition

const (
	ADMIN Role = iota + 1
	USER
)

var roleName = map[Role]string{
	ADMIN: "超级管理员",
	USER:  "普通用户",
}

type Role uint

func (role *Role) String() string {
	return roleName[*role]
}
