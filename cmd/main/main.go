package main

import (
	"github.com/devashishRaj/norm"
	"log"
)

func main() {
	err := norm.NomadmetricsBulksend()
	if err != nil {
		log.Fatal(err)
	}
}