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
				Path:    "/",
				Value:   "Created as test cookie at " + time.Now().String(),
				Expires: time.Now().Add(24 * time.Hour),
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
