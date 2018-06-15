package common

/*
	MOve to 'package common'
*/

import (
	"bytes"
	"encoding/json"
)

// GetAsJsonF is a debug function to inspect objects, will return formatted JSON
func GetAsJsonF(pIFC interface{}) string {
	var labJSON []byte
	var lsText string
	// var lErr error

	labJSON, _ = json.MarshalIndent(pIFC, "", "    ")

	lsText = bytes.NewBuffer(labJSON).String()
	// log.Println(lsText)
	return lsText
}

// GetAsJson is a debug function to inspect objects
func GetAsJson(pIFC interface{}) string {
	var labJSON []byte
	var lsText string
	// var lErr error

	labJSON, _ = json.Marshal(pIFC)

	lsText = bytes.NewBuffer(labJSON).String()
	// log.Println(lsText)
	return lsText
}
