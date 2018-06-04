package webServer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// get these vars from config
var miWebPort int
var msRedirectHost string

type custTransport struct{}

func (t *custTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	// Can request be modified?
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	// Can response be modified?
	return response, err
}

func StartProxy(piPort int, piWebPort int) {
	miWebPort = piWebPort
	var lRevProxyMux = http.NewServeMux()
	msRedirectHost = fmt.Sprintf("localhost:%d", miWebPort)

	lRevProxyMux.HandleFunc("/", proxyHandler)
	lRevProxySvr := &http.Server{
		Addr:    fmt.Sprintf(":%d", piPort),
		Handler: lRevProxyMux,
		// ReadTimeout:  650,
		// WriteTimeout: 650,
		// IdleTimeout:  650,
	}
	fmt.Println("Starting Reverse Proxy on port", fmt.Sprintf(":%d", piPort))
	lRevProxySvr.ListenAndServe()
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyHandler: Path =", r.URL.Path)

	fmt.Println("urlstring =", fmt.Sprintf("http://localhost:%d", miWebPort))
	urlstr := fmt.Sprintf("http://%s", msRedirectHost)
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Println("url parse error", err)
	}

	// Skip for now, till fileserver is enabled
	if r.URL.Path == "/favicon.ico" {
		fmt.Println("proxyHandler: Skipping call for /favicon.ico")
		return
	}

	/* Do this in custom Director
	fmt.Print("r.Host = { Before:", r.Host)
	r.Host = u.Host
	fmt.Println(", after: ", r.Host, "}")
	*/

	r.Header.Set("X-proxy-ReqVal", "Request value set by proxy")
	w.Header().Set("X-proxy-HeaderVal", "Header value set by proxy")

	revProxy := httputil.NewSingleHostReverseProxy(u)
	revProxy.Transport = &custTransport{}
	revProxy.Director = custHeader
	revProxy.ModifyResponse = setResponseHeader

	revProxy.ServeHTTP(w, r)

}

// custHeader will modify request object before forwarding to proxied server
func custHeader(req *http.Request) {
	/*
	   PLAN:
	       Check authentication
	       set payload object in header as JSON string
	       remove the request before sending back to client
	*/

	req.Header.Add("X-Origin-Host", req.Host)

	req.Header.Set("Host", msRedirectHost)
	req.URL.Scheme = "http"
	req.URL.Host = msRedirectHost

	req.Header.Add("X-Forwarded-Host", req.URL.Host)
	req.Header.Add("X-custHeader-time", "{todo:{gettime:pending,formattime:pending}}")
}

func setResponseHeader(w *http.Response) error {
	// TODO:
	w.Header.Add("X-Proxy-Add-srh", "setResponseHeader")
	return nil
}
