package common

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

// RecoverFromPanic if called as a 'defer' function
// before running potential panicky code,
// will gracefully log the error and prevent ap crash
func RecoverFromPanic() {
	log.Println("common.RecoverFromPanic: Enter", time.Now().String())

	// Get error causing panic
	lErrPanic := recover()
	if lErrPanic != nil {
		fmt.Println("Recovering from error")
	}

	// Write Error message
	// including link to safe default page
	log.Println("common.RecoverFromPanic: Exit", time.Now().String())
}

func WriteErrorResponse(pStackTrace []byte) {
	log.Println("common.WriteErrorResponse: Enter", time.Now().String())
	fmt.Println(bytes.NewBuffer(pStackTrace).String())
	log.Println("common.WriteErrorResponse: Exit", time.Now().String())
}
