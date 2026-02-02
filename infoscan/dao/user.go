package dao

// IUserDAO 用户数据访问接口
type IUserDAO interface {
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(username string) (*User, error)
	// CreateUser 创建新用户
	CreateUser(user *User) error
	// GetAllUsers 获取所有用户
	GetAllUsers() ([]*User, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *User) error
	// DeleteUser 删除用户
	DeleteUser(id uint) error
}
