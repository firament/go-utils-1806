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
	ClientIP     string
	ClientHost   string
	ConnProtocol string
	ConnScheme   string
	Headers      []entryLine
	Cookies      []entryLine
	QueryString  []entryLine
	RawDumpLines []string
}
type entryLine struct {
	Name  string
	Value string
}

// DumpRequestOutData is an experimental feature to inspect the Outbound Request snapshot
func DumpRequestOutData(pReqOut *http.Request) (rsOutput string) {
	labRespDump, err := httputil.DumpRequestOut(pReqOut, false)
	if err != nil {
		fmt.Printf("DumpRequestOutData: Dump error is %+v \n", err)
	}
	lsRespString := bytes.NewBuffer(labRespDump).String()
	rsOutput = utils1806.GetAsJson(strings.Split(lsRespString, "\r\n"))
	// fmt.Println("Response as String = ", lsRespString)
	return
}

// DumpResponseData is an experimental feature to inspect the Response snapshot
func DumpResponseData(pResp *http.Response) (rsOutput string) {
	labRespDump, err := httputil.DumpResponse(pResp, true)
	if err != nil {
		fmt.Printf("DumpResponseData: Dump error is %+v \n", err)
	}
	lsRespString := bytes.NewBuffer(labRespDump).String()
	rsOutput = utils1806.GetAsJson(strings.Split(lsRespString, "\r\n"))
	// fmt.Println("Response as String = ", lsRespString)
	return
}

// DumpRequestData will return Request contents in readable format
func DumpRequestData(r *http.Request) (rsOutput string) {

	log.Println("DumpRequestData - Begin")
	var lOutput requestData

	lOutput.ClientIP = r.RemoteAddr + " | " + r.Header.Get("X-Forwarded-For")
	lOutput.ClientHost = r.Host
	lOutput.ConnProtocol = r.Proto
	// lOutput.ConnScheme = bytes.NewBuffer(r.TLS.TLSUnique).String()

	lOutput.Headers = getHeaders(r.Header)
	lOutput.Cookies = getCookies(r.Cookies())
	lOutput.QueryString = getQueryParms(r.URL.RawQuery)

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
	var lsCleanOP string

	log.Printf("Processing %d Cookies", len(pCookies))
	rCookies = make([]entryLine, len(pCookies))
	var lEntry entryLine

	for liIndex, lCookie = range pCookies {
		lsCleanOP = strings.Replace(utils1806.GetAsJson(lCookie), "\n", "", -1)
		lsCleanOP = strings.Replace(lsCleanOP, "\"", "'", -1)
		lEntry = entryLine{
			Name:  lCookie.Name,
			Value: lsCleanOP,
		}
		rCookies[liIndex] = lEntry

		// fmt.Printf("Entry %d: Key = %s, Val = %s\n", liIndex, lEntry.Name, lEntry.Value)
		fmt.Printf("Raw Value is %+v \n", lCookie)
	}
	return
}

func getQueryParms(pQryString string) (rQueryParms []entryLine) {
	var lasQryParms []string
	var lsKV string
	var lasKVPair []string
	var liIndex int

	if len(strings.Trim(pQryString, " ")) == 0 {
		log.Printf("No Query Parms to process.")
		return
	}
	lasQryParms = strings.Split(pQryString, "&")

	log.Printf("Processing %d Query Parms", len(lasQryParms))
	rQueryParms = make([]entryLine, len(lasQryParms))
	var lEntry entryLine

	for liIndex, lsKV = range lasQryParms {
		lasKVPair = strings.Split(lsKV, "=")
		if len(lasKVPair) < 1 {
			lasKVPair[1] = ""
		}
		lEntry = entryLine{
			Name:  lasKVPair[0],
			Value: lasKVPair[1],
		}
		rQueryParms[liIndex] = lEntry

		// fmt.Printf("Entry %d: Key = %s, Val = %s\n", liIndex, lEntry.Name, lEntry.Value)

	}
	return
}
