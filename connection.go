package sutando

import (
	"crypto/tls"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection interface {
	DSN(cfg *tls.Config) (string, bool) /* string: dsn, bool: isPem */
	Database() string
	SetupOption(*options.ClientOptions)
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

	OptionHandler func(*options.ClientOptions)
}

func (c Conn) DSN(cfg *tls.Config) (string, bool) {
	suffix := ""
	prefix := "mongodb:"
	if c.Srv {
		prefix = "mongodb+srv:"
	}
	dsn := fmt.Sprintf("%s//%s:%s@%s:%d/%s",
		prefix,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)

	if c.Srv || c.Port == 0 {
		dsn = fmt.Sprintf("%s//%s:%s@%s/%s",
			prefix,
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

	return dsn + suffix, pem
}

func (c Conn) Database() string {
	return c.DB
}

func (c Conn) SetupOption(opt *options.ClientOptions) {
	if c.OptionHandler == nil {
		return
	}
	c.OptionHandler(opt)
}
