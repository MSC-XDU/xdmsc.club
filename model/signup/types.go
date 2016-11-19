package signup

import "time"

// LogIn 保存了用户的登录信息。
type Register struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Id       UserToken `json:"id"`
}

// 用于访问用户其他数据的索引
type UserToken uint64

// Profile 保存用户的基本个人信息
type Profile struct {
	Email          string    `json:"email"`
	EmailValidated bool      `json:"-"`
	PhoneNumber    string    `json:"phoneNumber"`
	PhoneValidated bool      `json:"-"`
	Name           string    `json:"name"`
	Sex            bool      `json:"sex"`
	Age            int       `json:"age"`
	HomeTown       string    `json:"homeTown"`
	Nation         string    `json:"nation"`
	Birthday       time.Time `json:"birthday, string"`
	QQ             string    `json:"qq"`
	Major          string    `json:"major"`
	UpdatedAt      time.Time `json:"updatedTime, omitempty"`
}

// SignUp 保存报名表信息
type SignUp struct {
	Department   string    `json:"department"`
	Introduction string    `json:"introduction"`
	Skills       string    `json:"skills"`
	Attachments  []string  `json:"attachments, omitempty"`
	Character    string    `json:"character"`
	UpdateAt     time.Time `json:"updateTime"`
	Idea         string    `json:"idea"`
	Achievement  string    `json:"achievement"`
}
