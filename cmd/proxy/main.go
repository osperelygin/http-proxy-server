package main

import (
	"http-proxy-server/internal/app/proxy"
	"log"
)

func main() {
	if err := proxy.Start(); err != nil {
		log.Fatalln(err)
	}
}
