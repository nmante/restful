package main

import (
	"log"
)

func main() {
	a, err := NewApp()

	if err != nil {
		log.Fatalf("Could not build Application: %s", err.Error())
	}

	a.Configure()
	a.Start()
}
