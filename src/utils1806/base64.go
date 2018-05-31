package utils1806

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

type b64Status struct {
	File       string
	SizeBinary int64
	SizeBase64 int64
	ErrData    string
	Base64Text string
}

// GetBase64 returns string
func GetBase64(w http.ResponseWriter, r *http.Request) {
	// Code is too short to write seperate functions
	switch lsMethod := r.Method; lsMethod {
	case "GET":
		http.Error(w, fmt.Sprintf("Feature %s is not yet implemented!", r.URL), http.StatusNotImplemented)
		return // redundant, but good coding practice
	case "POST":
		// Get uploaded file
		var lErr error
		var lsErrStr string
		var lStatus b64Status
		var lfpostedFile multipart.File
		var lfPostFileHeader *multipart.FileHeader
		var liBytes int
		var labData []byte
		var lsB64Text string

		lErr = r.ParseMultipartForm(5242880)
		if lErr != nil {
			lStatus.ErrData = fmt.Sprintf("FAIL Parsing Form, with error %+v\n", lErr)
			lsErrStr = getAsJson(lStatus)
			log.Println(lsErrStr)
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}

		lfpostedFile, lfPostFileHeader, lErr = r.FormFile("BinaryFile")
		if lErr != nil {
			lStatus.ErrData = fmt.Sprintf("FAIL Opening file, with error %+v\n", lErr)
			lsErrStr = getAsJson(lStatus)
			log.Println(lsErrStr)
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}
		log.Printf("Processing File %s containing %d bytes.\n", lfPostFileHeader.Filename, lfPostFileHeader.Size)
		lStatus.File = lfPostFileHeader.Filename
		lStatus.SizeBinary = lfPostFileHeader.Size

		// read content into array
		labData = make([]byte, lStatus.SizeBinary)

		liBytes, lErr = lfpostedFile.Read(labData)
		if lErr != nil {
			lStatus.ErrData = fmt.Sprintf("FAIL Reading file, with error %+v\n", lErr)
			lsErrStr = getAsJson(lStatus)
			log.Println(lsErrStr)
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}
		if int64(liBytes) != lStatus.SizeBinary {
			lStatus.ErrData = fmt.Sprintf("FAIL Size mismatch on reading file, expected %d got %d", liBytes, lStatus.SizeBinary)
			lsErrStr = getAsJson(lStatus)
			log.Println(lsErrStr)
			http.Error(w, lsErrStr, http.StatusExpectationFailed)
			return
		}

		lsB64Text = base64.StdEncoding.EncodeToString(labData)
		lStatus.Base64Text = lsB64Text
		lStatus.SizeBase64 = int64(len(lsB64Text))

		// Legacy Code

		// Write results back to caller
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Disposition", "inline")
		// w.Header().Add("Status-Code", strconv.Itoa(http.StatusMethodNotAllowed))
		w.Write([]byte(getAsJson(lStatus)))
		w.WriteHeader(http.StatusOK)

	default: // we do not allow any method as default, if not coded, it is an error
		http.Error(w, fmt.Sprintf("Request Method not allowed. Check documentation for valid methods!"), http.StatusMethodNotAllowed)
	}
}

// getAsJson is a debug function to inspect objects
func getAsJson(pIFC interface{}) string {
	var labJSON []byte
	var lsText string
	// var lErr error

	labJSON, _ = json.MarshalIndent(pIFC, "", "    ")

	lsText = bytes.NewBuffer(labJSON).String()
	// log.Println(lsText)
	return lsText
}
