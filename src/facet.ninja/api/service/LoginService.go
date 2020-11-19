package service

// taken from https://medium.com/wesionary-team/jwt-authentication-in-golang-with-gin-63dbc0816d55
type LoginService interface {
	LoginUser(email string, password string) bool
}
type loginInformation struct {
	email    string
	password string
}

func StaticLoginService() LoginService {
	return &loginInformation{
		email:    "bikash.dulal@wesionary.team",
		password: "testing",
	}
}
func (info *loginInformation) LoginUser(email string, password string) bool {
	return info.email == email && info.password == password
}