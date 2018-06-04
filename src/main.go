package main

import (
	"log"
	"utils1806/webServer"
)

var miWebPort int
var miProxyPort int

func main() {

	// Accept public port from args
	if miProxyPort < 8000 {
		miProxyPort = 9090
	}
	miWebPort = miProxyPort + 5

	webServer.StartWebServer(miWebPort)
	webServer.StartProxy(miProxyPort, miWebPort)
	log.Println("BYE! Stopping execution.")
}
