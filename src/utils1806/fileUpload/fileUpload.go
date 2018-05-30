package fileUpload

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

const CPass string = "PASS"
const CFAIL string = "FAIL"

// csvStatus will be used only in this file, for now.
type csvStatus struct {
	Status  string // PASS, FAIL or any other
	Row     int
	Col     int
	ErrData string
	// Error   error
}

// bulkUploadStatus will be used only in this file, for now.
type bulkUploadStatus struct {
	File      string
	Size      int64
	Rows      int
	RowStatus []csvStatus
}

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
		var lsErrStr string
		var lfpostedFile multipart.File
		var lfPostFileHeader *multipart.FileHeader
		var lStatus bulkUploadStatus
		var lRowStatus csvStatus

		lErr = r.ParseMultipartForm(5242880)
		if lErr != nil {
			lsErrStr = fmt.Sprintf("FAIL Parsing Form, with error %+v\n", lErr)
			lRowStatus.Row = -1
			lRowStatus.Status = lsErrStr
			lStatus.RowStatus = append(lStatus.RowStatus, lRowStatus)
			lsErrStr = getStatusAsStr(lStatus, true, "uploadCSV")
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}

		lfpostedFile, lfPostFileHeader, lErr = r.FormFile("ApiDataFileCsv")
		if lErr != nil {
			lsErrStr = fmt.Sprintf("FAIL Opening file, with error %+v\n", lErr)
			lRowStatus.Row = -1
			lRowStatus.Status = lsErrStr
			lStatus.RowStatus = append(lStatus.RowStatus, lRowStatus)
			lsErrStr = getStatusAsStr(lStatus, true, "uploadCSV")
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}
		log.Printf("Processing File %s containing %d bytes.\n", lfPostFileHeader.Filename, lfPostFileHeader.Size)
		lStatus.File = lfPostFileHeader.Filename
		lStatus.Size = lfPostFileHeader.Size

		var csvReader *csv.Reader
		csvReader = csv.NewReader(lfpostedFile)
		csvReader.ReuseRecord = false

		vsResult := uploadCSVPost(csvReader, lStatus)

		// Write results back to caller
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Disposition", "inline")
		w.Write([]byte(vsResult))
		w.WriteHeader(http.StatusOK)

	default: // we do not allow any method as default, if not coded, it is an error
		http.Error(w, fmt.Sprintf("Request Method not allowed. Check documentation for valid methods!"), http.StatusMethodNotAllowed)
	}
}

func uploadCSVGet(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Feature %s is not yet implemented!", r.URL), http.StatusNotImplemented)
	return // redundant, but good coding practice
}

func uploadCSVPost(pCsvReader *csv.Reader, pStatus bulkUploadStatus) string {
	var lasRow []string
	var lErr error
	var lsErrStr string
	var liRowCounter int
	// var lStatus bulkUploadStatus
	var lRowStatus csvStatus

	// Get Header row
	lasRow, lErr = pCsvReader.Read() // discard for now, need to validate column defs
	if lasRow == nil || lErr != nil {
		// Bad condition
		lsErrStr = fmt.Sprintf("FAIL Reading file header, with error %+v", lErr)
		lRowStatus.Row = -1
		lRowStatus.Status = lsErrStr
		pStatus.RowStatus = append(pStatus.RowStatus, lRowStatus)
		log.Println(lsErrStr)
		return getStatusAsStr(pStatus, true, "uploadCSVPost")

	}

	// process rows
	for {
		lasRow, lErr = pCsvReader.Read()
		if lasRow == nil {
			log.Printf("Processed %d records. No more records to process.\n", liRowCounter)
			pStatus.Rows = liRowCounter
			break
		} // stop when there are no more rows
		liRowCounter++

		// check for ​appropriate error
		if lErr != nil {
			lsErrStr = fmt.Sprintf("FAIL Reading file, with error %+v", lErr)
			lRowStatus.Row = -1
			lRowStatus.Status = lsErrStr
			pStatus.RowStatus = append(pStatus.RowStatus, lRowStatus)
			log.Println(lsErrStr)
			return getStatusAsStr(pStatus, true, "uploadCSVPost")
		}

		// process the record
		// 'lasRow' now contains data on one API
		lRowStatus = addAPIFromCsv(lasRow)
		lRowStatus.Row = liRowCounter
		pStatus.RowStatus = append(pStatus.RowStatus, lRowStatus)

		/* // NOT needed, will be done in called function
		if liRowStatus == http.StatusOK && lRowStatus.Status == CPass {
			log.Printf("PASS Processed Record %d\n", liRowCounter)
		} else {
			log.Printf("FAIL Record %d, Got Status code = %d, with status = %+v.\n ", liRowCounter, liRowStatus, lRowStatus)
			log.Println("Proceeding to next record.")
		}
		*/
	}

	// Consolidate logs and return
	return getStatusAsStr(pStatus, true, "uploadCSVPost")
}

func addAPIFromCsv(psAPIData []string) (rStatus csvStatus) {
	// parse the row and add API to system, as PUBLIC only

	log.Println("TODO: Process data.")

	//psAPIData[0] will be raw JSON

	// Validate all columns for size
	// pStatus.Col = 99, set status of error column

	// build and Add x-info object to JSON
	log.Println("Company", psAPIData[1])
	/*
		log.Println("x-App-Name", psAPIData[2])
		log.Println("x-App-URL", psAPIData[3])
		log.Println("x-App-Description", psAPIData[4])
		log.Println("x-API-Name", psAPIData[5])
		log.Println("x-API-Documentation-URL", psAPIData[6])
		log.Println("x-Company-Github", psAPIData[7])
		log.Println("x-additionalHostNames", psAPIData[8])
	*/
	// rStatus.Error = errors.New(fmt.Sprintf("Company = %s", psAPIData[1]))
	rStatus.ErrData = fmt.Sprintf("x-Company = %s", psAPIData[1])

	// Use transaction - begin

	// Insert into API table
	// Insert into INFO table

	// Use transaction - end

	// Insert into Mongo, if txn sucessful

	// TODO: use custom struct to return status

	// Set all-is-well flags
	rStatus.Status = CPass
	log.Println(getAsJson(rStatus))
	return
}

func getStatusAsStr(pStatus bulkUploadStatus, pbLogToDBAlso bool, pSender string) (rStatusText string) {
	var lsJsonText []byte
	var lErr error

	lsJsonText, lErr = json.MarshalIndent(pStatus, "", "    ")
	if lErr != nil {
		rStatusText = ""
	} else {
		rStatusText = bytes.NewBuffer(lsJsonText).String()
	}

	if pbLogToDBAlso {
		//utilities.WriteLogEntryDB(pSender, rStatusText)
	}
	return
}

// getAsJson is a debug function to inspect objects
func getAsJson(pIFC interface{}) string {
	var labJSON []byte
	var lsText string
	// var lErr error

	labJSON, _ = json.MarshalIndent(pIFC, "", "    ")

	lsText = bytes.NewBuffer(labJSON).String()
	log.Println(lsText)
	return lsText
}
