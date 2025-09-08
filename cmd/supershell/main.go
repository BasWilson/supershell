package main

import (
	"log"
	"os"

	"supershell/internal/cli"
)

func main() {
	if err := cli.Execute(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
