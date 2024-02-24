package main

import (
	"http-proxy-server/internal/app/webapi"
	"log"
)

func main() {
	if err := webapi.Start(); err != nil {
		log.Fatalln(err)
	}
}
