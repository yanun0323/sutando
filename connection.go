package sutando

import (
	"crypto/tls"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection provides a definition of connect configuration.
type Connection interface {
	DSN(cfg *tls.Config) (dns string, isPem bool)
	Database() string
	SetupClientOptions(*options.ClientOptions)
}

// Conn is an implementation of the Connection interface
// for connecting mongo db through the host and port.
type Conn struct {
	Username  string
	Password  string
	Host      string
	Port      uint   /* leave blank if you already add port in host */
	DB        string /* database name */
	Pem       string /* optional */
	AdminAuth bool

	ClientOptionsHandler func(*options.ClientOptions)
}

func (c Conn) DSN(cfg *tls.Config) (string, bool) {
	host := c.Host
	if c.Port != 0 {
		host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	}

	sf, pem := suffix(c.AdminAuth, cfg, c.Pem)
	return fmt.Sprintf("mongodb://%s:%s@%s/%s%s",
		c.Username,
		c.Password,
		host,
		c.DB,
		sf,
	), pem
}

func suffix(adminAuth bool, cfg *tls.Config, pem string) (string, bool) {
	if cfg.RootCAs.AppendCertsFromPEM([]byte(pem)) {
		return "?ssl=true&replicaSet=rs0&readpreference=secondaryPreferred", true
	}

	if adminAuth {
		return "?authSource=admin", false
	}

	return "", false
}

func (c Conn) Database() string {
	return c.DB
}

func (c Conn) SetupClientOptions(opt *options.ClientOptions) {
	if c.ClientOptionsHandler != nil {
		c.ClientOptionsHandler(opt)
	}
}

// ConnSrv is an implementation of the Connection interface
// for connecting mongo db through the SRV.
type ConnSrv struct {
	Username  string
	Password  string
	Host      string
	DB        string /* database name */
	Pem       string /* optional */
	AdminAuth bool

	ClientOptionsHandler func(*options.ClientOptions)
}

func (c ConnSrv) DSN(cfg *tls.Config) (string, bool) {
	sf, pem := suffix(c.AdminAuth, cfg, c.Pem)
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s%s",
		c.Username,
		c.Password,
		c.Host,
		c.DB,
		sf,
	), pem
}

func (c ConnSrv) Database() string {
	return c.DB
}

func (c ConnSrv) SetupClientOptions(opt *options.ClientOptions) {
	if c.ClientOptionsHandler == nil {
		return
	}
	c.ClientOptionsHandler(opt)
}
