package sutando

import (
	"crypto/tls"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection interface {
	DSN(cfg *tls.Config) (string, bool) /* string: dsn, bool: isPem */
	Database() string
	SetupClientOptions(*options.ClientOptions)
}

type Conn struct {
	Username  string
	Password  string
	Host      string
	Port      uint /* leave blank if you already add port in host or using SRV*/
	DB        string
	Pem       string /* optional */
	AdminAuth bool
	Srv       bool

	ClientOptionsHandler func(*options.ClientOptions)
}

func (c Conn) DSN(cfg *tls.Config) (string, bool) {
	prefix := "mongodb:"
	suffix := ""
	if c.Srv {
		prefix = "mongodb+srv:"
	}

	if c.AdminAuth {
		suffix = "?authSource=admin"
	}

	pem := cfg.RootCAs.AppendCertsFromPEM([]byte(c.Pem))
	if pem {
		suffix = "?ssl=true&replicaSet=rs0&readpreference=secondaryPreferred"
	}

	dsn := fmt.Sprintf("%s//%s:%s@%s:%d/%s%s",
		prefix,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
		suffix,
	)

	if c.Srv || c.Port == 0 {
		dsn = fmt.Sprintf("%s//%s:%s@%s/%s%s",
			prefix,
			c.Username,
			c.Password,
			c.Host,
			c.DB,
			suffix,
		)
	}

	return dsn, pem
}

func (c Conn) Database() string {
	return c.DB
}

func (c Conn) SetupClientOptions(opt *options.ClientOptions) {
	if c.ClientOptionsHandler == nil {
		return
	}
	c.ClientOptionsHandler(opt)
}
