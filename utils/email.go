package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	ReplyTo string   `json:"reply_to,omitempty"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type resendError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

// SendEmail sends an HTML email via the Resend API.
// The smtpPassword parameter is kept for backwards-compatibility but is unused;
// the Resend API key is read from the RESEND_API_KEY environment variable.
func SendEmail(fromAddress, toAddress, subject, body, replyToAddress, _ string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable is not set")
	}

	log.Printf("Preparing to send email from %s to %s with subject: %s", fromAddress, toAddress, subject)

	payload := resendRequest{
		From:    fromAddress,
		To:      []string{toAddress},
		Subject: subject,
		HTML:    body,
	}
	if replyToAddress != "" {
		payload.ReplyTo = replyToAddress
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Resend API: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var resendErr resendError
		if jsonErr := json.Unmarshal(respBody, &resendErr); jsonErr == nil && resendErr.Message != "" {
			return fmt.Errorf("Resend API error %d: %s", resendErr.StatusCode, resendErr.Message)
		}
		return fmt.Errorf("Resend API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Email sent successfully to %s", toAddress)
	return nil
}
