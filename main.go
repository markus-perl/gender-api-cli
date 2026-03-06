package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const baseURL = "https://gender-api.com/v2"

// GenderResponse represents the standard response from the /gender endpoint
type GenderResponse struct {
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	FullName    string  `json:"full_name,omitempty"`
	Email       string  `json:"email,omitempty"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	Details     struct {
		Duration string `json:"duration"`
	} `json:"details"`
}

// CountryOfOriginResponse represents the response from the /country-of-origin endpoint
type CountryOfOriginResponse struct {
	FirstName       string  `json:"first_name"`
	Gender          string  `json:"gender"`
	Probability     float64 `json:"probability"`
	CountryOfOrigin []struct {
		CountryName string  `json:"country_name"`
		Country     string  `json:"country"`
		Probability float64 `json:"probability"`
	} `json:"country_of_origin"`
	Details struct {
		Duration string `json:"duration"`
	} `json:"details"`
}

// Version will be injected at compile time
var Version = "dev"

// StatsResponse represents the response from the /statistic endpoint
type StatsResponse struct {
	IsLimitReached   bool `json:"is_limit_reached"`
	RemainingCredits int  `json:"remaining_credits"`
	UsageLastMonth   struct {
		Date        string `json:"date"`
		CreditsUsed int    `json:"credits_used"`
	} `json:"usage_last_month"`
}

func main() {
	var (
		firstName string
		fullName  string
		email     string
		country   string
		locale    string
		ip        string
		outFormat string
		flagKey   string
		origin    bool
		stats     bool
		version   bool
	)

	flag.StringVar(&firstName, "first_name", "", "First name to query")
	flag.StringVar(&fullName, "full_name", "", "Full name to query")
	flag.StringVar(&email, "email", "", "Email address to query")
	flag.StringVar(&country, "country", "", "ISO 3166 ALPHA-2 Country Code")
	flag.StringVar(&locale, "locale", "", "Browser Locale")
	flag.StringVar(&ip, "ip", "", "IP address for localization")
	flag.StringVar(&outFormat, "out", "text", "Output format: text or json")
	flag.StringVar(&flagKey, "key", "", "Gender-API authorization token")
	flag.BoolVar(&origin, "origin", false, "Query country of origin endpoint instead of standard gender endpoint")
	flag.BoolVar(&stats, "stats", false, "Query account statistics")
	flag.BoolVar(&version, "version", false, "Print the version of the CLI tool")
	flag.Parse()

	if version {
		fmt.Printf("Gender-API CLI Client v%s\n", Version)
		return
	}

	apiKey := getAPIKey(flagKey)
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not find Gender-API key.")
		fmt.Fprintln(os.Stderr, "Please provide it via -key flag, GENDER_API_KEY environment variable, or ~/.gender-api-key file.")
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	if stats {
		handleStats(client, apiKey, outFormat)
		return
	}

	if origin {
		handleOrigin(client, apiKey, firstName, fullName, email, outFormat)
		return
	}

	if firstName == "" && fullName == "" && email == "" {
		fmt.Fprintln(os.Stderr, "Error: You must provide either -first_name, -full_name, or -email.")
		os.Exit(1)
	}

	handleGender(client, apiKey, firstName, fullName, email, country, locale, ip, outFormat)
}

func doRequest(client *http.Client, apiKey, method, url string, payload map[string]string) ([]byte, error) {
	var bodyReader io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func handleGender(client *http.Client, apiKey, firstName, fullName, email, country, locale, ip, outFormat string) {
	payload := make(map[string]string)
	if firstName != "" {
		payload["first_name"] = firstName
	} else if fullName != "" {
		payload["full_name"] = fullName
	} else if email != "" {
		payload["email"] = email
	}

	if country != "" {
		payload["country"] = country
	}
	if locale != "" {
		payload["locale"] = locale
	}
	if ip != "" {
		payload["ip"] = ip
	}

	body, err := doRequest(client, apiKey, "POST", baseURL+"/gender", payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if outFormat == "json" {
		fmt.Println(string(body))
		return
	}

	var res GenderResponse
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	fmt.Println("--- Result ---")
	if res.FirstName != "" {
		fmt.Printf("First Name: %s\n", res.FirstName)
	}
	if res.LastName != "" {
		fmt.Printf("Last Name: %s\n", res.LastName)
	}
	if res.FullName != "" {
		fmt.Printf("Full Name: %s\n", res.FullName)
	}
	if res.Email != "" {
		fmt.Printf("Email: %s\n", res.Email)
	}
	fmt.Printf("Gender: %s\n", res.Gender)
	fmt.Printf("Accuracy: %.0f%%\n", res.Probability*100)
	fmt.Printf("Duration: %s\n", res.Details.Duration)
}

func handleOrigin(client *http.Client, apiKey, firstName, fullName, email, outFormat string) {
	if firstName == "" && fullName == "" && email == "" {
		fmt.Fprintln(os.Stderr, "Error: You must provide -first_name, -full_name, or -email for origin queries.")
		os.Exit(1)
	}

	payload := make(map[string]string)
	if firstName != "" {
		payload["first_name"] = firstName
	}
	if fullName != "" {
		payload["full_name"] = fullName
	}
	if email != "" {
		payload["email"] = email
	}

	body, err := doRequest(client, apiKey, "POST", baseURL+"/country-of-origin", payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if outFormat == "json" {
		fmt.Println(string(body))
		return
	}

	var res CountryOfOriginResponse
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	fmt.Println("--- Result ---")
	fmt.Printf("First Name: %s\n", res.FirstName)
	fmt.Printf("Gender: %s (%.0f%%)\n", res.Gender, res.Probability*100)

	fmt.Println("Countries of Origin:")
	for _, c := range res.CountryOfOrigin {
		fmt.Printf("  - %s (%s): %.0f%%\n", c.CountryName, c.Country, c.Probability*100)
	}
	fmt.Printf("Duration: %s\n", res.Details.Duration)
}

func handleStats(client *http.Client, apiKey, outFormat string) {
	body, err := doRequest(client, apiKey, "GET", baseURL+"/statistic", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if outFormat == "json" {
		fmt.Println(string(body))
		return
	}

	var res StatsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	fmt.Println("--- Account Statistics ---")
	fmt.Printf("Limit Reached: %v\n", res.IsLimitReached)
	fmt.Printf("Remaining Credits: %d\n", res.RemainingCredits)
	fmt.Printf("Last Month Usage (%s): %d credits\n", res.UsageLastMonth.Date, res.UsageLastMonth.CreditsUsed)
}

func getAPIKey(flagKey string) string {
	// 1. Command Line Flag
	if flagKey != "" {
		return flagKey
	}

	// 2. Environment Variable
	envKey := os.Getenv("GENDER_API_KEY")
	if envKey != "" {
		return envKey
	}

	// 3. Config File
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".gender-api-key")
		data, err := os.ReadFile(configPath)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return ""
}
