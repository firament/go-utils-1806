package ReqDebug

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"utils1806"
)

type requestData struct {
	Headers      []entryLine
	Cookies      []entryLine
	QueryString  []entryLine
	RawDumpLines []string
}
type entryLine struct {
	Name  string
	Value string
}

// DumpRequestData will return in readable format
func DumpRequestData(r *http.Request) (rsOutput string) {

	log.Println("DumpRequestData - Begin")
	var lOutput requestData
	lOutput.Headers = getHeaders(r.Header)
	lOutput.Cookies = getCookies(r.Cookies())
	lOutput.QueryString = getQueryParms(strings.Split(r.URL.RawQuery, "&"))

	// Get a copy, we dont want to consume the body
	reqCopy, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Printf("DumpRequestData: Dump error is %+v \n", err)
	}
	lOutput.RawDumpLines = strings.Split(bytes.NewBuffer(reqCopy).String(), "\r\n")

	rsOutput = utils1806.GetAsJson(lOutput)
	// rsOutput += fmt.Sprintf("\nBody Dump is \n%+v", bytes.NewBuffer(reqCopy).String())
	return
}

func getHeaders(pHeaders http.Header) (rHeaders []entryLine) {
	var lsKey string
	var lasVals []string
	var liIndex int

	log.Printf("Processing %d Headers", len(pHeaders))
	rHeaders = make([]entryLine, len(pHeaders))
	var lEntry entryLine

	for lsKey, lasVals = range pHeaders {
		lEntry = entryLine{
			Name:  lsKey,
			Value: strings.Join(lasVals, ", "),
		}
		rHeaders[liIndex] = lEntry
		liIndex++

		// fmt.Printf("Entry %d: Key = %s, Val = %s", liIndex, lEntry.Name, lEntry.Value)
	}
	return
}

func getCookies(pCookies []*http.Cookie) (rCookies []entryLine) {
	var lCookie *http.Cookie
	var liIndex int

	log.Printf("Processing %d Cookies", len(pCookies))
	rCookies = make([]entryLine, len(pCookies))
	var lEntry entryLine

	for liIndex, lCookie = range pCookies {
		lEntry = entryLine{
			Name:  lCookie.Name,
			Value: lCookie.Value,
		}
		rCookies[liIndex] = lEntry

		fmt.Printf("Entry %d: Key = %s, Val = %s\n", liIndex, lEntry.Name, lEntry.Value)
		fmt.Printf("Raw Value is %+v \n", lCookie)
	}
	return
}

func getQueryParms(pQryParms []string) (rQueryParms []entryLine) {
	var lsKV string
	var lasKVPair []string
	var liIndex int

	log.Printf("Processing %d Query Parms", len(pQryParms))
	rQueryParms = make([]entryLine, len(pQryParms))
	var lEntry entryLine

	for liIndex, lsKV = range pQryParms {
		lasKVPair = strings.Split(lsKV, "=")
		lEntry = entryLine{
			Name:  lasKVPair[0],
			Value: lasKVPair[1],
		}
		rQueryParms[liIndex] = lEntry

		// fmt.Printf("Entry %d: Key = %s, Val = %s\n", liIndex, lEntry.Name, lEntry.Value)

	}
	return
}
