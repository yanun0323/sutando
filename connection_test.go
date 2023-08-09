package sutando

import (
	"crypto/tls"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConn(t *testing.T) {
	testCases := []struct {
		desc             string
		conn             Connection
		expectedDSN      string
		expectedIsPem    bool
		expectedDatabase string
	}{
		{
			desc: "testDSN_Simple",
			conn: Conn{
				Username:  "username",
				Password:  "password",
				Host:      "localhost",
				Port:      27017,
				DB:        "test",
				Pem:       "",
				AdminAuth: false,
				Srv:       false,
			},
			expectedDSN:      "mongodb://username:password@localhost:27017/test",
			expectedIsPem:    false,
			expectedDatabase: "test",
		},
		{
			desc: "testDSN_AdminAuth",
			conn: Conn{
				Username:  "username",
				Password:  "password",
				Host:      "localhost",
				Port:      27017,
				DB:        "test",
				Pem:       "",
				AdminAuth: true,
				Srv:       false,
			},
			expectedDSN:      "mongodb://username:password@localhost:27017/test?authSource=admin",
			expectedIsPem:    false,
			expectedDatabase: "test",
		},
		{
			desc: "testDSN_Pem",
			conn: Conn{
				Username:  "username",
				Password:  "password",
				Host:      "localhost",
				Port:      27017,
				DB:        "test",
				Pem:       _pem,
				AdminAuth: false,
				Srv:       false,
			},
			expectedDSN:      "mongodb://username:password@localhost:27017/test?ssl=true&replicaSet=rs0&readpreference=secondaryPreferred",
			expectedIsPem:    true,
			expectedDatabase: "test",
		},
		{
			desc: "testDSN_SRV",
			conn: Conn{
				Username:  "username",
				Password:  "password",
				Host:      "url",
				Port:      27017,
				DB:        "test",
				Pem:       "",
				AdminAuth: true,
				Srv:       true,
			},
			expectedDSN:      "mongodb+srv://username:password@url/test?authSource=admin",
			expectedIsPem:    false,
			expectedDatabase: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Log(tc.desc)
			cfg := &tls.Config{
				RootCAs: x509.NewCertPool(),
			}
			dsn, isPem := tc.conn.DSN(cfg)
			assert.Equal(t, tc.expectedDSN, dsn)
			assert.Equal(t, tc.expectedIsPem, isPem)
			assert.Equal(t, tc.expectedDatabase, tc.conn.Database())
		})
	}
}

const _pem = `
-----BEGIN CERTIFICATE-----
MIIEBzCCAu+gAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgZQxCzAJBgNVBAYTAlVT
MRAwDgYDVQQHDAdTZWF0dGxlMRMwEQYDVQQIDApXYXNoaW5ndG9uMSIwIAYDVQQK
DBlBbWF6b24gV2ViIFNlcnZpY2VzLCBJbmMuMRMwEQYDVQQLDApBbWF6b24gUkRT
MSUwIwYDVQQDDBxBbWF6b24gUkRTIGFwLWVhc3QtMSBSb290IENBMB4XDTE5MDIx
NzAyNDcwMFoXDTIyMDYwMTEyMDAwMFowgY8xCzAJBgNVBAYTAlVTMRMwEQYDVQQI
DApXYXNoaW5ndG9uMRAwDgYDVQQHDAdTZWF0dGxlMSIwIAYDVQQKDBlBbWF6b24g
V2ViIFNlcnZpY2VzLCBJbmMuMRMwEQYDVQQLDApBbWF6b24gUkRTMSAwHgYDVQQD
DBdBbWF6b24gUkRTIGFwLWVhc3QtMSBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAOcJAUofyJuBuPr5ISHi/Ha5ed8h3eGdzn4MBp6rytPOg9NVGRQs
O93fNGCIKsUT6gPuk+1f1ncMTV8Y0Fdf4aqGWme+Khm3ZOP3V1IiGnVq0U2xiOmn
SQ4Q7LoeQC4lC6zpoCHVJyDjZ4pAknQQfsXb77Togdt/tK5ahev0D+Q3gCwAoBoO
DHKJ6t820qPi63AeGbJrsfNjLKiXlFPDUj4BGir4dUzjEeH7/hx37na1XG/3EcxP
399cT5k7sY/CR9kctMlUyEEUNQOmhi/ly1Lgtihm3QfjL6K9aGLFNwX35Bkh9aL2
F058u+n8DP/dPeKUAcJKiQZUmzuen5n57x8CAwEAAaNmMGQwDgYDVR0PAQH/BAQD
AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFFlqgF4FQlb9yP6c+Q3E
O3tXv+zOMB8GA1UdIwQYMBaAFK9T6sY/PBZVbnHcNcQXf58P4OuPMA0GCSqGSIb3
DQEBCwUAA4IBAQDeXiS3v1z4jWAo1UvVyKDeHjtrtEH1Rida1eOXauFuEQa5tuOk
E53Os4haZCW4mOlKjigWs4LN+uLIAe1aFXGo92nGIqyJISHJ1L+bopx/JmIbHMCZ
0lTNJfR12yBma5VQy7vzeFku/SisKwX0Lov1oHD4MVhJoHbUJYkmAjxorcIHORvh
I3Vj5XrgDWtLDPL8/Id/roul/L+WX5ir+PGScKBfQIIN2lWdZoqdsx8YWqhm/ikL
C6qNieSwcvWL7C03ri0DefTQMY54r5wP33QU5hJ71JoaZI3YTeT0Nf+NRL4hM++w
Q0veeNzBQXg1f/JxfeA39IDIX1kiCf71tGlT
-----END CERTIFICATE-----
`
