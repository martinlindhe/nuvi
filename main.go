package main

import (
	"flag"
	"log"

	"github.com/martinlindhe/nuvi/news"
)

var (
	verboseFlag = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	if err := news.ImportPublishedDocuments(*verboseFlag); err != nil {
		log.Fatal(err)
	}
}
