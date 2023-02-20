package assessment

import (
	"context"
	"fmt"
	"log"
	"os"

	recaptcha "cloud.google.com/go/recaptchaenterprise/apiv1"
	recaptchapb "google.golang.org/genproto/googleapis/cloud/recaptchaenterprise/v1"
)

/** PerformAssesment
 * Create an assessment to analyze the risk of an UI action.
 *
 * @param projectID: GCloud Project ID
 * @param recaptchaSiteKey: Site key obtained by registering a domain/app to use recaptcha services.
 * @param token: The token obtained from the client on passing the recaptchaSiteKey.
 */
func PerformAssesment(projectID string, recaptchaSiteKey string, token string) (float32, error) {

	// Create the recaptcha client.
	// TODO: To avoid memory issues, move this client generation outside
	// of this example, and cache it (recommended) or call client.close()
	// before exiting this method.
	ctx := context.Background()
	client, err := recaptcha.NewClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("Error creating reCAPTCHA client : %w", err)
	}
	defer client.Close()

	// Set the properties of the event to be tracked.
	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: recaptchaSiteKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	// Build the assessment request.
	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", projectID),
	}

	response, err := client.CreateAssessment(
		ctx,
		request)

	if err != nil {
		return 0, fmt.Errorf("Error creating assessment : %w", err)
	}

	// Check if the token is valid.
	if response.TokenProperties.Valid == false {
		return 0, fmt.Errorf("Assessment failed: %v",
			response.TokenProperties.InvalidReason)
	}

	// Get the risk score and the reason(s).
	// For more information on interpreting the assessment,
	// see: https://cloud.google.com/recaptcha-enterprise/docs/interpret-assessment
	log.Printf("The reCAPTCHA score for this token is:  %v",
		response.RiskAnalysis.Score)

	for _, reason := range response.RiskAnalysis.Reasons {
		log.Printf("%s", reason.String())
	}

	return response.RiskAnalysis.Score, nil

}

func Assess(token string) (float32, error) {
	projectId := os.Getenv("RC_PROJECT_ID")
	siteKey := os.Getenv("RC_SITEKEY")
	return PerformAssesment(projectId, siteKey, token)
}
