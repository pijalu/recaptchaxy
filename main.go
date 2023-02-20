package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pijalu/recaptchaxy/proxy"
)

func main() {
	listen := os.Getenv("RC_LISTEN")
	proxy := proxy.Proxy{
		Target: os.Getenv("RC_TARGET"),
	}

	log.Fatal(http.ListenAndServe(listen, &proxy))
}
