package main

import (
	metadataService "audigo/services/metadata"
	"flag"
	"fmt"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 8000, "Port for metadata service")
	flag.Parse()
	fmt.Printf("Starting API at http://localhost:%d\n", port)
	metadataService.Run(port)
}
