package webServer

import (
	"fmt"
	"net/http"
)

func StartWebServer(piPort int) {
	if piPort < 8000 {
		piPort = 8040
	}
	websrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", piPort),
		Handler: getAllEndpoints(),
		// ReadTimeout:  650,
		// WriteTimeout: 650,
		// IdleTimeout:  650,
	}
	fmt.Println("Starting web server on port", fmt.Sprintf(":%d", piPort))
	go websrv.ListenAndServe()
}

/*
w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

cfg := &tls.Config{
  MinVersion: tls.VersionTLS12,
}
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
}

go run `go env GOROOT`/src/crypto/tls/generate_cert.go --host localhost


Generate private key (.key)
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048
# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
sopenssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

Path to CA certs on Ubuntu
/etc/ssl/certs/ca-certificates.crt

https://blog.charmes.net/post/reverse-proxy-go/
https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702

*/

func getAllEndpoints() *http.ServeMux {
	var lSrvMux *http.ServeMux
	lSrvMux = http.NewServeMux()

	addTestEndPoints(lSrvMux)
	addFileUploadEndPoints(lSrvMux)

	return lSrvMux
}

// setMethodError will post log messages on the exception raised by the application
// the header tags are expected to be handled by the client appropiately
func setMethodError(w http.ResponseWriter, r *http.Request) (rErrData string) {
	// Dump URL and query parms
	// Dump all headers
	// Dump Cookies
	// Dump User data

	// Write error message
	// These header messages work
	w.Header().Add("x-app-err-mesg", fmt.Sprintf("Method invoked '%s' is not supported", r.Method))
	w.Header().Add("x-app-redirect", fmt.Sprintf("/landing-page.html"))
	w.WriteHeader(http.StatusMethodNotAllowed)

	var lsErrStr string
	lsErrStr = fmt.Sprintf("setMethodError: Called by Endpoint '%s' with Method '%s'", r.URL.Path, r.Method)
	w.Write([]byte(lsErrStr))

	// return the JSON formatted error dump, in case the caller want to consume it
	return
}
