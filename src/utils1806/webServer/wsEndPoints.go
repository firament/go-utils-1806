package webServer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"utils1806/ReqDebug"
)

func addEndPoints(pSrvMux *http.ServeMux) {

	// webserver.Handle("/", http.FileServer(http.Dir(httpContentBasePath))
	lsCWD, _ := os.Getwd()
	log.Println(lsCWD)

	pSrvMux.Handle("/", http.FileServer(http.Dir(httpContentBasePath)))

	// Add function teapot, a simple test to check life
	pSrvMux.HandleFunc("/teapot", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		var lsText string
		lsText = "I am a little teapot... time now is " + time.Now().String()
		w.Write([]byte(lsText))
	})

	// Add function SwitchTest, to test fall through
	pSrvMux.HandleFunc("/switchtest", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		r.ParseForm()
		var lsTag string = r.Form.Get("tagval")
		switch lsTag {
		case "a":
			fallthrough
		case "b":
			fallthrough
		case "c":
			fmt.Println("cases A-C")
		case "d":
			fallthrough
		case "e":
			fmt.Println("case D, E")
		default:
			fmt.Println("other than A-E")
		}

	})

	// Add function doPanic
	pSrvMux.HandleFunc("/dopanic", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		// Create an unhandled panic, regardless of method
		// log.Panicln("Raising a PANIC at", time.Now().String())	// causes non-fatal panic
		panic("Raising a PANIC at" + time.Now().String()) // causes non-fatal panic
		println("Survived Panic")
		r.Response.Header.Get("utils1806-redir-to") // causes fatal panic? not from here
	})

	// Add function doRedir
	pSrvMux.HandleFunc("/doredir", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		http.Redirect(w, r, "reqdebug", http.StatusTemporaryRedirect)
		fmt.Println("Inspect response in proxy atfet this")
	})

	// Add function ReqDebug
	pSrvMux.HandleFunc("/reqdebug", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		switch r.Method {
		case http.MethodGet:
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Content-Disposition", "inline")
			w.Write([]byte(ReqDebug.DumpRequestData(r)))
			w.WriteHeader(http.StatusOK)
		default:
			setMethodError(w, r)
		}
	})

	// Add function cookieadd
	pSrvMux.HandleFunc("/cookieadd", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		switch r.Method {
		case http.MethodGet:
			http.SetCookie(w, &http.Cookie{
				Name:    fmt.Sprintf("test-cki-%d", time.Now().Second()),
				Path:    msCookiePath,
				Domain:  msCookieDomain,
				Value:   "Created as test cookie at " + time.Now().String(),
				Expires: time.Now().Add(2 * time.Minute),
			})
			// Repeat cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "test-repeat",
				Path:    msCookiePath,
				Domain:  msCookieDomain,
				Value:   "Created as test cookie at " + time.Now().String(),
				Expires: time.Now().Add(2 * time.Minute),
			})
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Add("Content-Disposition", "inline")
			w.Write([]byte(fmt.Sprintf("Done, add test cookie. at %s", time.Now().String())))
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			w.Write([]byte("TO Be Done, Try again some other time."))
			w.WriteHeader(http.StatusNotImplemented)
		default:
			setMethodError(w, r)
		}
	})

}
