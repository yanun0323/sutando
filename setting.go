package sutando

import (
	"crypto/tls"
	"fmt"
)

type Connection interface {
	DSN(cfg *tls.Config) (string, bool) /* string: dsn, bool: isPem */
	Database() string
}

type Conn struct {
	Username  string
	Password  string
	Host      string
	Port      uint
	DB        string
	Pem       string
	AdminAuth bool
}

func (c Conn) DSN(cfg *tls.Config) (string, bool) {
	suffix := ""
	prefix := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)

	if c.Port == 0 {
		prefix = fmt.Sprintf("mongodb://%s:%s@%s/%s",
			c.Username,
			c.Password,
			c.Host,
			c.DB,
		)
	}

	if c.AdminAuth {
		suffix = "?authSource=admin"
	}

	var pem bool
	if pem = cfg.RootCAs.AppendCertsFromPEM([]byte(c.Pem)); pem {
		suffix = "?ssl=true&replicaSet=rs0&readpreference=secondaryPreferred"
	}

	return prefix + suffix, pem
}

func (c Conn) Database() string {
	return c.DB
}
