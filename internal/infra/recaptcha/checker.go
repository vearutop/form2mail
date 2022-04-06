package recaptcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bool64/stats"
)

type siteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type Config struct {
	SiteKey    string  `split_words:"true"`
	SecretKey  string  `split_words:"true"`
	V3         bool    `split_words:"true"`
	V3MinScore float64 `split_words:"true" default:"0.5"`
}

type Checker interface {
	CheckToken(ctx context.Context, responseToken string) error
}

// V2V3Checker checks recaptcha status.
type V2V3Checker struct {
	Transport http.RoundTripper
	Config    Config
	Stats     stats.Tracker
}

// RecaptchaChecker is a service provider.
func (c *V2V3Checker) RecaptchaChecker() Checker {
	return c
}

func (c *V2V3Checker) CheckToken(ctx context.Context, responseToken string) error {
	if c.Config.SecretKey == "" {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.google.com/recaptcha/api/siteverify", nil)
	if err != nil {
		return err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", c.Config.SecretKey)
	q.Add("response", responseToken)
	req.URL.RawQuery = q.Encode()

	tr := c.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}

	// Make request
	resp, err := tr.RoundTrip(req)
	if err != nil {
		if c.Stats != nil {
			c.Stats.Add(ctx, "recaptcha_failed", 1, "reason", "request")
		}

		return err
	}
	defer resp.Body.Close()

	// Decode response.
	var body siteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		if c.Stats != nil {
			c.Stats.Add(ctx, "recaptcha_failed", 1, "reason", "decode")
		}
		return err
	}

	// Check recaptcha verification success.
	if !body.Success {
		if c.Stats != nil {
			c.Stats.Add(ctx, "recaptcha_failed", 1, "reason", "unsuccessful")
		}
		return errors.New("unsuccessful recaptcha verify request")
	}

	if c.Config.V3 {
		minScore := c.Config.V3MinScore
		if minScore == 0 {
			minScore = 0.5
		}

		if c.Stats != nil {
			c.Stats.Add(ctx, "recaptcha_score_total", body.Score)
			c.Stats.Add(ctx, "recaptcha_score_count", 1)
		}

		// Check response score.
		if body.Score < minScore {
			if c.Stats != nil {
				c.Stats.Add(ctx, "recaptcha_failed", 1, "reason", "low_score")
			}

			return fmt.Errorf("lower received score than expected: %.2f", body.Score)
		}
	}

	return nil
}
