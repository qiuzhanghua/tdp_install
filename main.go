package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var destination = ""
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homeDir)

	flag.Usage = func() {
		fmt.Printf("Install TDP to Destination Folder\n\nUSAGE:\n%s [OPTIONS]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.StringVarP(&destination, "exdir", "d", homeDir, "Directory where TDP installed")
	flag.Parse()

	if len(destination) < 1 {
		destination = homeDir
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()
		if strings.HasPrefix(fileName, "tdp-") && strings.HasSuffix(fileName, ".zip") {
			// TODO unzip it to destination/tdp
			// TODO and others
		}
		fmt.Println(f.Name())
	}
}
