package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func fetchDataHandler(w http.ResponseWriter, r *http.Request) {

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenUrl := os.Getenv("TOKEN_URL")
	serviceUrl := os.Getenv("SERVICE_URL")

	log.Println("CLIENT_ID:", clientId)
	log.Println("CLIENT_SECRET:", clientSecret)
	log.Println("TOKEN_URL:", tokenUrl)
	log.Println("SERVICE_URL:", serviceUrl)

	queryParams := r.URL.Query()
	country := queryParams.Get("country")
	encodedCountry := url.QueryEscape(country)
	url := fmt.Sprintf("%s%s%s%s", serviceUrl, "/v3/covid-19/countries/", encodedCountry, "?strict=true")

	// Get the token
	token := getToken(clientId, clientSecret, tokenUrl)
	if token == nil {
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}
	var tokenResponse AccessTokenResponse
	err := json.Unmarshal([]byte(*token), &tokenResponse)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	accessToken := tokenResponse.AccessToken
	log.Println("Access Token:", accessToken)

	// Make a GET request to the REST service
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making GET request: %s", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %s", err), http.StatusInternalServerError)
		return
	}

	// Write the result to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getToken(clientID string, clientSecret string, tokenEndpoint string) *string {
	// Create a base64-encoded string for the Authorization header
	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}
	log.Println("Response:", string(body))
	token := string(body)
	return &token
}

func main() {
	// Create a new HTTP server
	http.HandleFunc("/info", fetchDataHandler)

	// Start the server on port 8080
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
