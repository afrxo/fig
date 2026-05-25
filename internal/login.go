package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/afrxo/fig/auth"
)

var LoginCommand = Command{
	Name:    "login",
	Summary: "Authenticate with your Git host",
	Usage:   fmt.Sprintf("%s login", CLI),
	Run:     runLogin,
}

const clientID = "Ov23li2K3sj9kJVgtfDv"

func runLogin(flags CommandFlags, args []string) error {
	deviceCode, userCode, verificationURI, interval, err := requestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to start login: %w", err)
	}

	fmt.Println()
	fmt.Printf("! First, copy your one-time code: %s\n", userCode)
	fmt.Printf("- Press Enter to open github.com in your browser...")
	_, _ = fmt.Scanln()
	_ = OpenURL(verificationURI)
	fmt.Println("Waiting for authentication...")

	token, err := pollForToken(deviceCode, interval)
	if err != nil {
		return err
	}

	username, err := fetchGitHubUsername(token)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}

	if err := auth.Save(auth.Credentials{Username: username, Token: token}); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Printf("✓ Authentication complete.\n")
	fmt.Printf("✓ Logged in as %s\n", username)
	return nil
}

func requestDeviceCode() (deviceCode, userCode, verificationURI string, interval int, err error) {
	req, err := http.NewRequest("POST", "https://github.com/login/device/code", strings.NewReader(
		`{"client_id":"`+clientID+`","scope":"repo"}`,
	))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json") // ← critical, without this GitHub returns form-encoded text

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		DeviceCode      string `json:"device_code"`
		UserCode        string `json:"user_code"`
		VerificationURI string `json:"verification_uri"`
		Interval        int    `json:"interval"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	return result.DeviceCode, result.UserCode, result.VerificationURI, result.Interval, nil
}

func pollForToken(deviceCode string, interval int) (string, error) {
	if interval == 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
			strings.NewReader(`{
        "client_id": "`+clientID+`",
        "device_code": "`+deviceCode+`",
        "grant_type": "urn:ietf:params:oauth:grant-type:device_code"
    }`),
		)
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		var result struct {
			AccessToken string `json:"access_token"`
			Error       string `json:"error"`
			Interval    int    `json:"interval"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		switch result.Error {
		case "":
			return result.AccessToken, nil // success
		case "authorization_pending":
			continue // user hasn't clicked authorize yet
		case "slow_down":
			interval = result.Interval // GitHub asked us to back off
		case "expired_token":
			return "", fmt.Errorf("code expired, run login again")
		case "access_denied":
			return "", fmt.Errorf("authorization denied by user")
		default:
			return "", fmt.Errorf("unexpected error: %s", result.Error)
		}
	}
}

func fetchGitHubUsername(token string) (string, error) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Login string `json:"login"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Login, nil
}
