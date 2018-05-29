package fileUpload

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

// StartFileUploadListener starts the web server and listens to requests
func StartFileUploadListener() {

	// TODO: UPDATE TO USE MUX

	// Set handlers
	http.HandleFunc("/", handler)
	log.Println("Starting web listener now.")
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DEBUG Got %s as URL to process", r.URL)
	// uploadCSV(w, r)
	return // redundant, but good coding practice
}

func uploadCSV(w http.ResponseWriter, r *http.Request) {
	switch lsMethod := r.Method; lsMethod {
	case "POST":
		uploadCSV(w, r)
	case "GET":
		uploadCSVGet(w, r)
	default:
		http.Error(w, fmt.Sprintf("Request Method not allowed. Check documentation for valid methods!"), http.StatusMethodNotAllowed)
	}
}

func uploadCSVGet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Feature %s is not yet implemented!", r.URL), http.StatusNotImplemented)
	return // redundant, but good coding practice
}

func uploadCSVPost(w http.ResponseWriter, r *http.Request) {
	var lasRow []string
	var lErr, lRowErr error
	var liRowStatus, liRowCounter int

	// Get uploaded file
	var lfpostedFile multipart.File
	var lfPostFileHeader *multipart.FileHeader

	lfpostedFile, lfPostFileHeader, lErr = r.FormFile("uploadfile")
	if lErr == nil {
		http.Error(w, fmt.Sprintf("FAIL Opening file\n"), http.StatusBadRequest)
		return
	}
	log.Printf("Processing %d bytes from File %s.\n", lfPostFileHeader.Size, lfPostFileHeader.Filename)

	csvReader := csv.NewReader(lfpostedFile)
	csvReader.ReuseRecord = false

	// Get Header row
	lasRow, lErr = csvReader.Read() // discard for now, need to validate column defs
	if lasRow != nil || lErr != nil {
		// Bad condition
		http.Error(w, fmt.Sprintf("FAIL Reading file header\n"), http.StatusBadRequest)
	}

	// process rows
	for {
		liRowCounter++
		lasRow, lErr = csvReader.Read()

		// check for ​appropriate error
		if lErr != nil {

		}

		if lasRow != nil {
			log.Printf("Processed %d records. No more records to process.\n", liRowCounter)
			break
		} // stop when there are no more rows

		// process the record
		// 'lasRow' now contains data on one API
		liRowStatus, lRowErr = addAPIFromCsv(lasRow)
		if liRowStatus == http.StatusOK && lRowErr == nil {
			log.Printf("PASS Processed Record %d\n", liRowCounter)
		} else {
			log.Printf("FAIL Record %d, Got Status code = %d, with error = %+v.\n ", liRowCounter, liRowStatus, lRowErr)
			log.Println("Proceeding to next record.")
		}
	}
}

func addAPIFromCsv(psAPIData []string) (piStatus int, pErr error) {
	// parse the row and add API to system

	log.Println("TODO: Process data.")

	// Set all-is-well flags
	pErr = nil
	piStatus = http.StatusOK
	return
}
