package webServer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime/debug"
	"strings"
	"time"
	"utils1806/common"

	"github.com/Jeffail/gabs"
	"github.com/pborman/uuid"
)

// get these vars from config
var miWebPort int
var msRedirectHost string
var msCookieDomain string = "localhost"
var msCookiePath string = "/"

type custTransport struct{}

func (t *custTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	log.Println("custTransport.RoundTrip: Enter", time.Now().String())
	// Can request be modified?
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		// TODO: Write appropiate error message
		common.WriteErrorResponse(debug.Stack())
		request.Header.Add("RoundTrip", "Error GUUID")
		return nil, err
	}
	// Can response be modified?
	response.Header.Add("X-RoundTrip-A", "Header added in custTransport.RoundTrip")

	log.Println("custTransport.RoundTrip: Exit", time.Now().String())
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
	log.Printf("webServer.ProxyHandler: Path = %s ; urlstring = %s", r.URL.Path, fmt.Sprintf("http://localhost:%d", miWebPort))
	fmt.Println("webServer.ProxyHandler: Enter", time.Now().String())
	r.Header.Add("X-proxyHandler-Req", "Header in Request Add in webServer.proxyHandler")
	w.Header().Add("X-proxyHandler-Resp-Add", "Header in Response Add in webServer.proxyHandler")
	w.Header().Set("X-proxyHandler-Resp-Set", "Header in Response set in webServer.proxyHandler")

	// TEMP Inspection writes
	fmt.Println("r.Host =", r.Host)
	// fmt.Println("r.Response =", r.Response)
	// fmt.Println (" =", r.)
	// fmt.Println (" =", r.)

	urlstr := fmt.Sprintf("http://%s", msRedirectHost)
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Println("url parse error", err)
		return
	}

	// Log Sample - begin

	// SESSION ID - For sessions before login happens
	// Reuse Value if cookie already exists
	var lsSessID string = uuid.NewUUID().String() // Use LongCodePT()
	var lSessCki, lLoginCki *http.Cookie
	var lErrCki error

	// Get Http Session cookie
	lSessCki, lErrCki = r.Cookie("utils1806-session-id")
	if lErrCki == nil {
		fmt.Println("Existing Session ID = ", lsSessID)
		lsSessID = lSessCki.Value
	} else {
		fmt.Println("Creating new Session ID = ", lsSessID)
	}
	// Update cookie, to slide expiry
	http.SetCookie(w, &http.Cookie{
		Name:     "utils1806-session-id",
		Value:    lsSessID,
		Path:     msCookiePath,
		Domain:   msCookieDomain,
		HttpOnly: false,
		Expires:  time.Now().Add(2048 * time.Hour),
		Secure:   false,
	})

	// TODO: Use value from config
	// Expires:  time.Now().Add(2048 * time.Hour),

	var lsPath = strings.ToLower(r.URL.Path)
	switch {
	case strings.HasSuffix(lsPath, "css"):
	case strings.HasSuffix(lsPath, "js"):
	case strings.HasSuffix(lsPath, "jpeg"):
	case strings.HasSuffix(lsPath, "jpg"):
	case strings.HasSuffix(lsPath, "png"):
	case strings.HasSuffix(lsPath, "ico"):
		// Do not log for these cases
	default:
		fmt.Println("LOG ENTRY WILL BE WRITTEN.")
		// Write the log here, and only here

		var logRoot string = "LOG-OBJECT"
		var logEntryJSON = gabs.New()
		logEntryJSON.Set(lsSessID, logRoot, "SessionID")
		logEntryJSON.Set(r.URL.EscapedPath(), logRoot, "RequsetPath")
		logEntryJSON.Set(r.Header.Get("X-Forwarded-For"), logRoot, "ForwardedFor")
		logEntryJSON.Set(r.Header.Get("User-Agent"), logRoot, "UserAgent")
		logEntryJSON.Set(r.Header.Get("Referer"), logRoot, "Referer")

		// Get App Session cookie
		lLoginCki, lErrCki = r.Cookie("AppAuthCookie")
		if lErrCki == nil {
			fmt.Println(lLoginCki.Value)
			// TODO Write values from cookie
			logEntryJSON.Set("", logRoot, "UserName")
			logEntryJSON.Set("0", logRoot, "UserID")
			logEntryJSON.Set("0", logRoot, "OrgID")
			logEntryJSON.Set("0", logRoot, "RoleID")
		} else {
			// Keep blank entries, to ensure fields availaible for processing
			logEntryJSON.Set("", logRoot, "User", "UserName")
			logEntryJSON.Set("0", logRoot, "User", "UserID")
			logEntryJSON.Set("0", logRoot, "User", "OrgID")
			logEntryJSON.Set("0", logRoot, "User", "RoleID")
		}

		log.Println(logEntryJSON.StringIndent("", "  ")) // For debugging
		// log.Println(logEntryJSON.String()) // for PROD

	} // switch - end

	// log sample - end

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

	revProxy.ServeHTTP(w, r) // Do the actual reverse proxying...

	lsRTErr := r.Header.Get("RoundTrip")
	if lsRTErr != "" {
		log.Println("Round Trip Error Code = ", lsRTErr)
	}

	// Inspect status - for redirects
	var lsRedirTag string
	// lsRedirTag = r.Response.Header.Get("utils1806-redir-to") This will cause an PANIC here
	lsRedirTag = w.Header().Get("utils1806-redir-to")
	if lsRedirTag != "" {
		log.Println("Redirecting request to", lsRedirTag)
	} else {
		fmt.Println("No Redirect detected")
	}
	/*
		Manage Redirects internally
		* Update request with redirect path and query string
		* Copy cookies, if needed
		* call self proxyHandler(w, r)
	*/

	// This cookie has no effect
	http.SetCookie(w, &http.Cookie{
		Name:     "apptag2",
		Value:    "Post - Some Unique String, maybe raw guuid",
		Path:     msCookiePath,
		Domain:   msCookieDomain,
		HttpOnly: false,
		Expires:  time.Now().Add(5 * time.Minute),
		Secure:   false,
	})
	// fmt.Println("hook for debugging")
	log.Println("webServer.ProxyHandler: Exit", time.Now().String())
}

// proxyRequestMgr will modify request object before forwarding to proxied server
func proxyRequestMgr(req *http.Request) {
	log.Println("webServer.proxyRequestMgr: Enter", time.Now().String())
	req.Header.Add("X-proxyRequestMgr", "Header added in webServer.proxyRequestMgr")
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

	lTagCki, lErr := req.Cookie("utils1806-session-id")
	if lErr == nil {
		fmt.Printf("Value of 'utils1806-session-id' Cookie is %s\n", lTagCki.Value)
	}
	lTagCki, lErr = req.Cookie("apptagX")
	if lErr == nil {
		fmt.Printf("Value of TagX Cookie is %s", lTagCki.Value)
	}

	log.Println("webServer.proxyRequestMgr: Exit", time.Now().String())
}

// proxyResponseUpdate adds or modifies headers
// // and-or adds or modifies cookies
// to response before returning to client
func proxyResponseUpdate(w *http.Response) error {
	log.Println("webServer.proxyResponseUpdate: Enter", time.Now().String())

	// Capture Redirect condition
	var lsReloc string = w.Header.Get("Location")
	if lsReloc != "" {
		log.Println("REDIRECTING TO LOCATION:", lsReloc)
		w.Header.Add("utils1806-redir-to", lsReloc)
		// w.Request.Header.Add("utils1806-redir-to", lsReloc)
	} else {
		fmt.Println("NO RELOCATION INFO FOUND.")
	}

	// w.StatusCode() // This could be used to capture redirects

	w.Header.Add("X-setResponseHeader", "Header added in webServer.setResponseHeader")
	log.Println("webServer.proxyResponseUpdate: Exit", time.Now().String())
	return nil
}
