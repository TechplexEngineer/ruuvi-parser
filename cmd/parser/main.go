package main

import (
	"encoding/json"
	"github.com/techplexengineer/ruuvi-go-parser"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	//v5FromDocs := "0201061BFF99040505941A5BC7B1FFE0001C043867366F2497ED4DFAE75678"
	bath := "02010011FF990403651321C4F8013B00ED03E30B23"
	measurement, err := parser.Parse(bath)
	if err != nil {
		return err
	}
	indent, err := json.MarshalIndent(measurement, "", "    ")
	if err != nil {
		return err
	}
	log.Print(string(indent))
	return nil
}
