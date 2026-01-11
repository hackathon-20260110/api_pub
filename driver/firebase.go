package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseAuthClient *auth.Client

var ErrFirebaseNotInitialized = errors.New("firebase: not initialized")

// Token exchange request/response structures for Firebase REST API
type signInWithCustomTokenRequest struct {
	Token             string `json:"token"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type signInWithCustomTokenResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

type firebaseAPIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type serviceAccountCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

func NewFirebaseAuth() {
	creds := serviceAccountCredentials{
		Type:                    os.Getenv("FIREBASE_SA_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_SA_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_SA_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FIREBASE_SA_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FIREBASE_SA_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_SA_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_SA_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_SA_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FIREBASE_SA_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FIREBASE_SA_CLIENT_X509_CERT_URL"),
		UniverseDomain:          os.Getenv("FIREBASE_SA_UNIVERSE_DOMAIN"),
	}

	if creds.ProjectID == "" || creds.PrivateKey == "" || creds.ClientEmail == "" {
		log.Fatal("firebase: required environment variables (FIREBASE_SA_PROJECT_ID, FIREBASE_SA_PRIVATE_KEY, FIREBASE_SA_CLIENT_EMAIL) are not set")
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("firebase: failed to marshal credentials: %v", err)
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credJSON))
	if err != nil {
		log.Fatalf("firebase: failed to initialize app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("firebase: failed to get auth client: %v", err)
	}

	firebaseAuthClient = client
}

func FirebaseAuthClient() (*auth.Client, error) {
	if firebaseAuthClient == nil {
		return nil, ErrFirebaseNotInitialized
	}
	return firebaseAuthClient, nil
}

func VerifyIDToken(ctx context.Context, idToken string) (userID string, email string, err error) {
	client, err := FirebaseAuthClient()
	if err != nil {
		return "", "", err
	}

	decodedToken, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", "", fmt.Errorf("firebase: failed to verify ID token: %w", err)
	}

	userID = decodedToken.UID

	if e, ok := decodedToken.Claims["email"].(string); ok {
		email = e
	}

	return userID, email, nil
}

// CreateCustomToken generates a Firebase custom token for the specified user ID.
// This token can be used for testing and development purposes.
// Returns the custom token string or an error if generation fails.
func CreateCustomToken(ctx context.Context, userID string) (string, error) {
	client, err := FirebaseAuthClient()
	if err != nil {
		return "", err
	}

	token, err := client.CustomToken(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("firebase: failed to create custom token: %w", err)
	}

	return token, nil
}

// ExchangeCustomTokenForIDToken exchanges a Firebase custom token for an ID token
// using the Firebase Authentication REST API.
// This allows server-side code to obtain ID tokens that can be used as Bearer tokens.
//
// Parameters:
//   - ctx: Context for the HTTP request
//   - customToken: The Firebase custom token to exchange
//
// Returns:
//   - idToken: The Firebase ID token (can be used as Bearer token)
//   - expiresIn: Token expiration time in seconds (e.g., "3600")
//   - error: Any error that occurred during the exchange
func ExchangeCustomTokenForIDToken(ctx context.Context, customToken string) (idToken string, expiresIn string, err error) {
	// Get Firebase Web API Key from environment
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	if apiKey == "" {
		return "", "", errors.New("firebase: FIREBASE_WEB_API_KEY environment variable is not set")
	}

	// Construct the Firebase REST API endpoint
	endpoint := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s", apiKey)

	// Prepare request body
	requestBody := signInWithCustomTokenRequest{
		Token:             customToken,
		ReturnSecureToken: true,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", fmt.Errorf("firebase: failed to marshal request body: %w", err)
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", "", fmt.Errorf("firebase: failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute HTTP request with 10-second timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("firebase: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("firebase: failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errorResp firebaseAPIErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			// If we can't parse the error response, return raw body
			return "", "", fmt.Errorf("firebase: token exchange failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}
		return "", "", fmt.Errorf("firebase: token exchange failed: %s (status: %d)", errorResp.Error.Message, errorResp.Error.Code)
	}

	// Parse successful response
	var response signInWithCustomTokenResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", "", fmt.Errorf("firebase: failed to parse response: %w", err)
	}

	// Validate that we received an ID token
	if response.IDToken == "" {
		return "", "", errors.New("firebase: received empty ID token from Firebase API")
	}

	return response.IDToken, response.ExpiresIn, nil
}
