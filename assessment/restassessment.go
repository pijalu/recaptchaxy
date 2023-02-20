package assessment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

type EnterpriseRestAssessment struct {
	ProjectID string
	ApiKey    string
	MinScore  float64
}

func (ea *EnterpriseRestAssessment) PerformEnterpriseRestAssessment(siteKey string, action string, token string) (float64, error) {
	url := fmt.Sprintf("https://recaptchaenterprise.googleapis.com/v1/projects/%s/assessments?key=%s",
		ea.ProjectID,
		ea.ApiKey)

	requestBody := fmt.Sprintf(`
	{
		"event": {
		  "token": "%s",
		  "siteKey": "%s",
		  "expectedAction": "%s"
		}
	}`, token, siteKey, action)

	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return 0, fmt.Errorf("error creating reCAPTCHA client : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", dump)
		return 0, fmt.Errorf("error calling reCAPTCHA service")
	}

	var x map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading reCAPTCHA results : %w", err)
	}

	json.Unmarshal(body, &x)
	ra, ok := x["riskAnalysis"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("error converting reCAPTCHA risk analysis results")
	}
	score, ok := ra["score"].(float64)
	if !ok {
		return 0, fmt.Errorf("error converting reCAPTCHA risk analysis score")
	}

	return score, nil
}
