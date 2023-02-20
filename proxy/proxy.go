package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pijalu/recaptchaxy/assessment"
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

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func tofloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		if str == "" {
			log.Printf("Using default recaptcha value of 0.5")
			return 0.5
		}
		log.Fatalf("Could not parse recaptcha level %s", str)
	}
	return f
}

var rea = assessment.EnterpriseRestAssessment{
	ProjectID: os.Getenv("RC_PROJECT_ID"),
	ApiKey:    os.Getenv("RC_APIKEY"),
	MinScore:  tofloat(os.Getenv("RC_MIN_SCORE")),
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	targetURL := fmt.Sprintf("%s%s", p.Target, req.RequestURI)

	if req.Method != "OPTIONS" {
		score, err := rea.PerformEnterpriseRestAssessment(
			req.Header.Get("x-recaptcha-site"),
			req.Header.Get("x-recaptcha-action"),
			req.Header.Get("x-recaptcha-token"))

		if err != nil {
			handleError(w, err)
			return
		}

		if score < rea.MinScore {
			handleError(w, err)
			return
		}
	}

	newReq, err := http.NewRequest(req.Method, targetURL, req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// Copy headers
	for key := range req.Header {
		for _, value := range req.Header[key] {
			newReq.Header.Add(key, value)
		}
	}
	for _, key := range hopHeaders {
		newReq.Header.Del(key)
	}

	// X-Forwarded-For
	host := req.RemoteAddr
	if prior, ok := newReq.Header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	newReq.Header.Set("X-Forwarded-For", host)

	client := http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		handleError(w, err)
		return
	}
	defer resp.Body.Close()

	for k := range resp.Header {
		for _, s := range resp.Header[k] {
			w.Header().Add(k, s)
		}
	}

	// Enforce CORS
	if req.Method == "OPTIONS" && len(resp.Header["Access-Control-Allow-Origin"]) == 0 {
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
	}

	w.WriteHeader(resp.StatusCode)
	len, err := io.Copy(w, resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	log.Printf("Proxy: %s %s - len(%d)", req.Method, targetURL, len)

}
