package main

import (
	"log"
	"utils1806/fileUpload"
)

func main() {
	fileUpload.StartFileUploadListener()
	log.Println("BYE! Stopping execution.")
}
