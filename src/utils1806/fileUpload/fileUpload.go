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
	http.ListenAndServe(":8085", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	ldbgURL := r.URL
	fmt.Println("URL = ", ldbgURL)

	switch lsURL := r.URL.Path; lsURL {
	case "/uploadCSV":
		uploadCSV(w, r)
	default:
		fmt.Fprintf(w, "DEBUG Got '%s' as URL to process", r.URL)
	}
	return // redundant, but good coding practice
}

// bootstrap code, for testing - END

func uploadCSV(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/uploadCSV" {
		http.Error(w, fmt.Sprintf("Unknown Endpoint '%s', Check documentation!", r.URL), http.StatusInternalServerError)
		return
	}

	// TODO: Allow only user designated
	// http.Error(w, fmt.Sprintf("Unknown Endpoint asked '%s', Check documentation!", r.URL), http.StatusInternalServerError)

	switch lsMethod := r.Method; lsMethod {
	case "GET":
		uploadCSVGet(w, r)
	case "POST":
		// Get uploaded file
		var lErr error
		var lfpostedFile multipart.File
		var lfPostFileHeader *multipart.FileHeader

		lErr = r.ParseMultipartForm(5242880)
		if lErr != nil {
			http.Error(w, fmt.Sprintf("FAIL Parsing Form, with error %+v\n", lErr), http.StatusExpectationFailed)
			return
		}

		lfpostedFile, lfPostFileHeader, lErr = r.FormFile("ApiDataFileCsv")
		if lErr != nil {
			http.Error(w, fmt.Sprintf("FAIL Opening file, with error %+v\n", lErr), http.StatusExpectationFailed)
			return
		}
		log.Printf("Processing File %s containing %d bytes.\n", lfPostFileHeader.Filename, lfPostFileHeader.Size)

		var csvReader *csv.Reader
		csvReader = csv.NewReader(lfpostedFile)
		csvReader.ReuseRecord = false

		uploadCSVPost(csvReader)

	default: // we do not allow any method as default, if not coded, it is an error
		http.Error(w, fmt.Sprintf("Request Method not allowed. Check documentation for valid methods!"), http.StatusMethodNotAllowed)
	}
}

func uploadCSVGet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Feature %s is not yet implemented!", r.URL), http.StatusNotImplemented)
	return // redundant, but good coding practice
}

func uploadCSVPost(pCsvReader *csv.Reader) {
	var lasRow []string
	var lErr, lRowErr error
	var liRowStatus, liRowCounter int

	// Get Header row
	lasRow, lErr = pCsvReader.Read() // discard for now, need to validate column defs
	if lasRow == nil || lErr != nil {
		// Bad condition
		// http.Error(w, fmt.Sprintf("FAIL Reading file header\n"), http.StatusNoContent)
		log.Println(fmt.Sprintf("FAIL Reading file header\n"))
		return
	}

	// process rows
	for {
		lasRow, lErr = pCsvReader.Read()
		if lasRow == nil {
			log.Printf("Processed %d records. No more records to process.\n", liRowCounter)
			break
		} // stop when there are no more rows
		liRowCounter++

		// check for ​appropriate error
		if lErr != nil {
			// http.Error(w, fmt.Sprintf("FAIL Reading file row\n"), http.StatusNoContent)
			log.Println(fmt.Sprintf("FAIL Reading file row\n"))
			return
		}

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

	// Consolidate logs and return
}

func addAPIFromCsv(psAPIData []string) (piStatus int, pErr error) {
	// parse the row and add API to system, as PUBLIC only

	log.Println("TODO: Process data.")

	//psAPIData[0] will be raw JSON

	// Insert into Mongo

	// Insert into API table
	// Insert into INFO table

	// TODO: use custom struct to return status

	// Set all-is-well flags
	pErr = nil
	piStatus = http.StatusOK
	return
}
