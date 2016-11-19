package signup

type ErrorCode uint

const (
	UserExist = iota
	UserNotExist
	PasswordError
	UsernameInvalid
	PasswordTooShort
	InformationInvalid
)

type Error struct {
	Code     ErrorCode `json:"errorCode"`
	ErrorStr string    `json:"error"`
}

var (
	UserExistErr          = Error{UserExist, "用户名已存在，请重新设置"}
	UserNotExistErr       = Error{UserNotExist, "不存在的用户名，请检查后重试"}
	PasswordErrorErr      = Error{PasswordError, "密码错误，请检查后重试"}
	UsernameInvalidErr    = Error{UsernameInvalid, "用户名不符合要求，请重新设置"}
	PasswordTooShortErr   = Error{PasswordTooShort, "密码太短了，请重新设置"}
	InformationInvalidErr = Error{InformationInvalid, "信息不完整，请检查"}
)

func (err Error) Error() string {
	return err.ErrorStr
}
