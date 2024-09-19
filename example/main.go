package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hlubek/readercomp"
)

func main() {
	result, err := readercomp.FilesEqual(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
