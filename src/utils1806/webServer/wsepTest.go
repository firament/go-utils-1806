package webServer

import (
	"log"
	"net/http"
	"utils1806/ReqDebug"
)

func addTestEndPoints(pSrvMux *http.ServeMux) {

	// Add function "/"
	pSrvMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		w.Write([]byte("Default endpoint called."))
	})

	// Add function ReqDebug
	pSrvMux.HandleFunc("/reqdebug", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Disposition", "inline")
		w.Write([]byte(ReqDebug.DumpRequestData(r)))
		w.WriteHeader(http.StatusOK)

		// For now, dont enforce method

		/*
			switch r.Method {
			case "GET":
			case "POST":
			default:
				// setMethodError(w, r)
			}
		*/
	})

	// Add function "A"
	pSrvMux.HandleFunc("/testa", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.RequestURI)
		switch r.Method {
		case "GET":
			w.Write([]byte(testGetA()))
		case "POST":
			w.Write([]byte(testPostA()))
		default:
			setMethodError(w, r)
		}
	})

	// Add function "B"
	pSrvMux.HandleFunc("/testb", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %s", r.URL.Path, r.Method, r.URL.String())
		switch r.Method {
		case "POST":
			w.Write([]byte(testPostB()))
		default:
			setMethodError(w, r)
		}
	})

	// Add function "C"
	pSrvMux.HandleFunc("/testc", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint %s called with Method %s, Full URI = %+v", r.URL.Path, r.Method, r.URL.RequestURI())
		switch r.Method {
		case "GET":
			log.Println("Req Header X-custHeader-time =", r.Header.Get("X-custHeader-time"))
			w.Write([]byte(testGetC()))
		default:
			setMethodError(w, r)
		}
	})

}

/*
End point definitions kept here for demo
in prod code, there would be in appropiate packages
*/
func testGetA() string {
	return "You have called 'testGetA'"
}

func testPostA() string {
	return "You have called 'testPostA'"
}

func testPostB() string {
	return "You have called 'testPostB'"
}

func testGetC() string {
	return "You have called 'testGetC'"
}
