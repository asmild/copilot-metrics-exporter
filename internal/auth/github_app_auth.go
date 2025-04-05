package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// GitHubAppAuth implements AuthProvider for GitHub Apps
type GitHubAppAuth struct {
	appID          int64
	installationID int64
	privateKey     *rsa.PrivateKey
	token          string
	expiry         time.Time
	client         *http.Client
}

// Token response from GitHub API
type tokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewGitHubAppAuth creates a new GitHub App auth provider
func NewGitHubAppAuth(appID, installationID int64, privateKeyPEM []byte) (*GitHubAppAuth, error) {
	// Parse private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	auth := &GitHubAppAuth{
		appID:          appID,
		installationID: installationID,
		privateKey:     key,
		client:         &http.Client{},
	}

	return auth, nil
}

// GetToken returns a valid installation token, refreshing if needed
func (a *GitHubAppAuth) GetToken() (string, error) {
	if a.token != "" && time.Now().Add(5*time.Minute).Before(a.expiry) {
		return a.token, nil
	}

	return a.refreshToken()
}

// refreshToken gets a new installation token
func (a *GitHubAppAuth) refreshToken() (string, error) {
	jwt, err := a.generateJWT()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", a.installationID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to get installation token: %s", resp.Status)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	a.token = tokenResp.Token
	a.expiry = tokenResp.ExpiresAt

	return a.token, nil
}

// generateJWT creates a signed JWT for GitHub App authentication
func (a *GitHubAppAuth) generateJWT() (string, error) {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		Issuer:    fmt.Sprintf("%d", a.appID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
