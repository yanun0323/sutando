package sutando

type Connection struct {
	Username  string
	Password  string
	Host      string
	Port      uint
	Database  string
	Pem       string
	AdminAuth bool
}
