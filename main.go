package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Nealm03/finance_tracker/transactions"
)

func main() {
	filePath := flag.String("file-path", "/Users/nealmorris/Downloads/5017_08022026.csv", "the file to process")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("file-path flag is required")
		os.Exit(1)
	}
	ingestFile(*filePath)
}

func ingestFile(filePath string) {

	dirName := path.Dir(filePath)
	dir := os.DirFS(dirName)
	fName := strings.TrimPrefix(strings.Split(filePath, dirName)[1], "/")
	transactionsImporter, err := transactions.NewLLloydsImporter(fName, dir, false)
	if err != nil {
		fmt.Printf("failed to create ingester: %v\n", err)
		return
	}

	results, err := transactionsImporter.Import(context.Background(), fName)
	if err != nil {

		fmt.Printf("failed to ingest file: %v\n", err)
		return
	}

	fmt.Printf("ingested transactions: %+v\n", results)
}
