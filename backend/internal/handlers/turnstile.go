package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// verifyTurnstile validates a Cloudflare Turnstile client token.
// Returns nil on success. If secret is empty (dev mode), always succeeds.
func verifyTurnstile(ctx context.Context, secret, clientToken, remoteIP string) error {
	if secret == "" {
		return nil // dev mode: skip verification
	}

	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", clientToken)
	if remoteIP != "" {
		form.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("siteverify request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !result.Success {
		return fmt.Errorf("turnstile verification failed")
	}
	return nil
}
