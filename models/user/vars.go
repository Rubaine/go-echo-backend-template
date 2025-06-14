package user

import (
	"time"

	"github.com/fatih/structs"
)

const TOKEN_EXPIRATION = time.Minute * 30

type (
	User struct {
		ID           int64  `structs:"id"`
		Email        string `structs:"email"`
		Firstname    string `structs:"firstname"`
		Lastname     string `structs:"lastname"`
		Password     string `structs:"-"`
		RecoverToken string `structs:"-"`
		Enable       bool   `structs:"-"`
		Admin        bool   `structs:"admin"`
	}

	UserList []User
)

func (user User) ToSelfWebDetail() map[string]any {
	return structs.Map(user)
}

func (user User) ToWeb() map[string]any {
	m := structs.Map(user)
	delete(m, "email")
	delete(m, "admin")
	return m
}

func (userList UserList) ToWeb() []map[string]any {
	m := make([]map[string]any, 0)
	for _, u := range userList {
		m = append(m, u.ToWeb())
	}
	return m
}
