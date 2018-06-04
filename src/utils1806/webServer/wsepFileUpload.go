package webServer

import (
	"net/http"
	"utils1806/fileUpload"
)

func addFileUploadEndPoints(pSrvMux *http.ServeMux) {

	// Add function "uploadCSV"
	pSrvMux.HandleFunc("uploadCSV", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fileUpload.UploadCSVGet(w, r)
		case "POST":
			fileUpload.UploadCSVPost(w, r)
		default:
			setMethodError(w, r)
		}
	})
}
