package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Proxy struct {
	Target string
}

func handleError(w http.ResponseWriter, err error) {
	log.Printf("Error doing request: %v", err)
	errorStr := fmt.Sprintf("%v", err)
	w.WriteHeader(500)
	w.Write([]byte(errorStr))
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	targetURL := fmt.Sprintf("%s%s", p.Target, req.RequestURI)
	log.Printf("Proxy: %s", targetURL)

	newReq, err := http.NewRequest(req.Method, targetURL, req.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	newReq.Header.Add("X-Forwarded-For", req.RemoteAddr)
	for key := range req.Header {
		// skip encoding
		if key == "Accept-Encoding" {
			continue
		}
		for _, value := range req.Header[key] {
			newReq.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		handleError(w, err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	for k := range resp.Header {
		for _, s := range resp.Header[k] {
			w.Header().Add(k, s)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
