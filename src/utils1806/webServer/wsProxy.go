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
	fmt.Printf("custTransport.RoundTrip: RoundTrip for %+v\n", request.URL)
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
	var lsCurrFn string = "webServer.ProxyHandler: "

	log.Printf("\n\n%sNew Requesuest \nEnter @ %s", lsCurrFn, time.Now().String())
	fmt.Printf("%s Path = %s ; urlstring = %s/n", lsCurrFn, r.URL.Path, fmt.Sprintf("http://localhost:%d", miWebPort))

	// REQUEST - Add header entry, and test in downstream for propogation
	r.Header.Add("X-proxyHandler-Req", "Header in Request Add in webServer.proxyHandler")
	// RESPONSE WRITER - Add header entry, and test in downstream for propogation
	w.Header().Add("X-proxyHandler-Resp-Add", "Header in Response Add in webServer.proxyHandler")
	w.Header().Set("X-proxyHandler-Resp-Set", "Header in Response set in webServer.proxyHandler")

	// TEMP Inspection writes
	fmt.Println(lsCurrFn, "r.Host =", r.Host)
	// fmt.Println("r.Response =", r.Response)
	// fmt.Println (" =", r.)
	// fmt.Println (" =", r.)

	// Set proxy routing
	urlstr := fmt.Sprintf("http://%s", msRedirectHost)
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Println(lsCurrFn, "url parse error", err)
		return
	}

	/***************************************************************************
	**                                                                        **
	** Log Sample - BEGIN                                                     **
	**                                                                        **
	***************************************************************************/

	// SESSION ID - For sessions before login happens
	var lsSessID string = uuid.NewUUID().String()
	var lSessCki, lLoginCki *http.Cookie
	var lErrCki error
	var lbWriteLog, lbIsEndpoint bool

	// Test for existing Http Session cookie
	lSessCki, lErrCki = r.Cookie("utils1806-session-id")
	if lErrCki == nil {
		fmt.Println(lsCurrFn, "Existing Session ID = ", lsSessID)
		lsSessID = lSessCki.Value
	} else {
		fmt.Println(lsCurrFn, "New Session ID = ", lsSessID)
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

	// TODO: Use value from config -> Expires:  time.Now().Add(2048 * time.Hour),

	var lsPath = strings.ToLower(r.URL.Path)   // need consistent case for text-comparison
	if r.URL.Path == "/" || r.URL.Path == "" { // will never be ""
		lsPath = "/index.html"
		fmt.Printf("%s Default Path: Replacing %s with %s", lsCurrFn, r.URL.Path, lsPath)
	}

	switch {
	//** Do not log for these cases **//
	case strings.HasSuffix(lsPath, ".css"):
	case strings.HasSuffix(lsPath, ".js"):
	case strings.HasSuffix(lsPath, ".jpeg"):
	case strings.HasSuffix(lsPath, ".jpg"):
	case strings.HasSuffix(lsPath, ".png"):
	case strings.HasSuffix(lsPath, ".ico"):

	//** Log for these cases **//
	case strings.HasSuffix(lsPath, ".htm"):
		lbWriteLog = true
	case strings.HasSuffix(lsPath, ".html"):
		lbWriteLog = true
	default:
		lbIsEndpoint = true // Flag will send this request to backing server
		lbWriteLog = true
	} // switch - end

	if lbWriteLog {
		fmt.Printf("%s REQUEST LOG ENTRY WILL BE WRITTEN. for path %s \n", lsCurrFn, lsPath)
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

		// log.Println(logEntryJSON.StringIndent("", "  ")) // For debugging, formatted
		log.Println(logEntryJSON.String()) // for PROD

	}

	/***************************************************************************
	**                                                                        **
	** Log Sample - END                                                       **
	**                                                                        **
	***************************************************************************/

	// Test entries to inspect propogation
	r.Header.Set("X-proxy-ReqVal", "Request value set by proxy")
	w.Header().Set("X-proxy-HeaderVal", "Header value set by proxy")

	if !lbIsEndpoint {
		r.Header.Set("X-proxy-static-req", "true")

		//** Serve static files directly **//

		// fmt.Printf("%s Resource requested: %s\n", lsCurrFn, r.RequestURI)
		var lsFile string = httpContentBasePath
		if r.URL.Path == "" || r.URL.Path == "/" {
			lsFile += "/index.html"
		} else {
			lsFile += r.URL.Path
		}
		// fmt.Printf("%s Serving File: %s\n", lsCurrFn, lsFile)
		fmt.Printf("%s Resource requested: %s, Serving File: %s\n", lsCurrFn, r.RequestURI, lsFile)
		http.ServeFile(w, r, lsFile)

		// Can do a return safely here... for file serves
		log.Println(lsCurrFn, "Exit", time.Now().String())
		return

	}

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
		log.Println(lsCurrFn, "Round Trip Error Code = ", lsRTErr)
	}

	// Inspect status - for redirects
	var lsRedirTag string
	lsRedirTag = w.Header().Get("utils1806-redir-to")
	if lsRedirTag != "" {
		log.Println(lsCurrFn, "Redirecting request to", lsRedirTag)
	} else {
		fmt.Println(lsCurrFn, "No Redirect detected")
	}
	fmt.Printf("%s X-PRU-Action == %s\n", lsCurrFn, w.Header().Get("X-PRU-Action"))
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

	log.Println(lsCurrFn, "Exit", time.Now().String())
}

// proxyRequestMgr will modify request object before forwarding to proxied server
func proxyRequestMgr(req *http.Request) {
	var lsCurrFn string = "webServer.proxyRequestMgr: "
	log.Println(lsCurrFn, "Enter", time.Now().String())

	// REQUEST - Add header entry, and test in downstream for propogation
	req.Header.Add("X-proxyRequestMgr", "Header in Request Add in webServer.proxyRequestMgr")
	/*
	   PLAN:
	       Check authentication
	       set payload object in r.context as object
	       // remove the request before sending back to client - not reqd
	*/
	req.Header.Add("X-PRM-Origin-Host", req.Host) // test entry to inspect propogation
	req.Header.Set("Host", msRedirectHost)
	req.URL.Scheme = "http"
	req.URL.Host = msRedirectHost

	// r.Header.Set("X-proxy-static-req", "true") // set in 'webServer.ProxyHandler'
	if req.Header.Get("X-proxy-static-req") != "true" {
		req.URL.Path = strings.ToLower(req.URL.Path) // Convert to lowercase, to allow readable client side casing, only for end-points
	}

	req.Header.Add("X-PRM-Forwarded-Host", req.URL.Host)
	req.Header.Add("X-PRM-time", "{time-info:{gettime:"+time.Now().String()+",formattime:"+time.Now().Format("Mon 6th June 2018 18:24 +0530")+"}}")

	// Add auth payload to context here
	req.Header.Add("X-PRM-Authenticated", "false")

	lTagCki, lErr := req.Cookie("utils1806-session-id")
	if lErr == nil {
		fmt.Printf("%s Value of 'utils1806-session-id' Cookie is %s\n", lsCurrFn, lTagCki.Value)
	}
	lTagCki, lErr = req.Cookie("apptagX")
	if lErr == nil {
		fmt.Printf("%s Value of TagX Cookie is %s\n", lsCurrFn, lTagCki.Value)
	}

	log.Println(lsCurrFn, "Exit", time.Now().String())
}

// proxyResponseUpdate adds or modifies headers
// // and-or adds or modifies cookies
// to response before returning to client
func proxyResponseUpdate(w *http.Response) error {
	var lsCurrFn string = "webServer.proxyResponseUpdate: "
	log.Println(lsCurrFn, "Enter", time.Now().String())

	// RESPONSE WRITER - Add header entry, and test in downstream for propogation
	w.Header.Add("X-PRU-Resp-Add", "Header in Response Add in webServer.proxyResponseUpdate")
	w.Header.Set("X-PRU-Resp-Set", "Header in Response set in webServer.proxyResponseUpdate")

	// Generalized codes for managing custom errors
	const (
		lcOK         int = 200
		lcErrRedir   int = 300
		lcErrMissing int = 400
		lcErrError   int = 500
	)

	var lsAction, lsReLocation string
	var liGenRetCode int

	// Capture Redirect condition, from upstream
	fmt.Printf("%s w.StatusCode = %d, w.Status = %s\n", lsCurrFn, w.StatusCode, w.Status)
	lsReLocation = w.Header.Get("Location")
	if lsReLocation != "" {
		log.Println(lsCurrFn, "Header - REDIRECTING TO LOCATION:", lsReLocation)
		w.Header.Add("utils1806-redir-to", lsReLocation)
	} else {
		fmt.Println(lsCurrFn, "Header - NO RELOCATION INFO FOUND.")
	}

	liRetCode := w.StatusCode
	switch { // Needs to be refined per project standard
	case liRetCode < 200: // Intermediate codes
		liGenRetCode = lcOK
		lsAction = ""
	case liRetCode < 300: // Normal, no action
		liGenRetCode = lcOK
		lsAction = ""
	case liRetCode < 400: // Do nothing here
		liGenRetCode = lcErrRedir
		lsAction = w.Header.Get("Location")
	case liRetCode < 500: // do a redirect
		liGenRetCode = lcErrMissing
		lsAction = "Redirect to Missing Resource page"
	case liRetCode < 600:
		liGenRetCode = lcErrError
		lsAction = "Redirect to System Error page"
	}
	fmt.Printf("%s liGenRetCode == %d, lsAction == %s\n", lsCurrFn, liGenRetCode, lsAction)

	w.Header.Set("X-PRU-Action", lsAction)
	log.Println(lsCurrFn, "Exit", time.Now().String())

	return nil
}
