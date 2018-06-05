package webServer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
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

// StartProxy starts the proxy server and listens on the given port
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

// proxyHandler configures the proxy server
// and will manage all requests that come in
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ProxyHandler: Path = %s ; urlstring = %s", r.URL.Path, fmt.Sprintf("http://localhost:%d", miWebPort))

	// fmt.Println("ProxyHandler: Path =", r.URL.Path)
	// fmt.Println("urlstring =", fmt.Sprintf("http://localhost:%d", miWebPort))
	urlstr := fmt.Sprintf("http://%s", msRedirectHost)
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Println("url parse error", err)
		return
	}

	// Test entries to inspect propogation
	r.Header.Set("X-proxy-ReqVal", "Request value set by proxy")
	w.Header().Set("X-proxy-HeaderVal", "Header value set by proxy")

	// Configure Proxy
	revProxy := httputil.NewSingleHostReverseProxy(u)
	revProxy.Transport = &custTransport{}
	revProxy.Director = proxyRequestMgr
	revProxy.ModifyResponse = proxyResponseUpdate

	/* Is this required? or will original context have the object?
	// Get Auth Context
	authCtx := context.WithValue(r.Context(), "key-context", "value-context")
	revProxy.ServeHTTP(w, r.WithContext(authCtx))
	*/

	revProxy.ServeHTTP(w, r)

	/*
		Manage Redirects internally
		* Update request with redirect path and query string
		* Copy cookies, if needed
		* call self proxyHandler(w, r)
	*/
	log.Println("hook for debugging")
}

// proxyRequestMgr will modify request object before forwarding to proxied server
func proxyRequestMgr(req *http.Request) {
	/*
	   PLAN:
	       Check authentication
	       set payload object in r.context as object
	       // remove the request before sending back to client - not reqd
	*/

	req.Header.Add("X-phm-Origin-Host", req.Host) // test entry to inspect propogation
	req.Header.Set("Host", msRedirectHost)
	req.URL.Scheme = "http"
	req.URL.Host = msRedirectHost
	req.URL.Path = strings.ToLower(req.URL.Path) // Convert to lowercase, to allow readable client side casing

	req.Header.Add("X-phm-Forwarded-Host", req.URL.Host)
	req.Header.Add("X-phm-time", "{time-info:{gettime:"+time.Now().String()+",formattime:"+time.Now().Format("Mon 6th June 2018 18:24 +0530")+"}}")

	// Add auth payload to context here
	req.Header.Add("X-phm-Authenticated", "false")

}

// proxyResponseUpdate adds or modifies headers
// // and-or adds or modifies cookies
// to response before returning to client
func proxyResponseUpdate(w *http.Response) error {
	// TODO:
	w.Header.Add("X-prh-Add", "setResponseHeader")
	return nil
}
