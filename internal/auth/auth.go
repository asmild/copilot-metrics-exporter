package auth

import (
	"errors"
)

// AuthProvider defines interface for authentication methods
type AuthProvider interface {
	Token() (string, error)
}

// NewAuthProvider creates appropriate auth provider based on config
func NewAuthProvider(config *config.Config) (AuthProvider, error) {
	if config.PAT != "" {
		return NewPATAuth(config.PAT), nil
	}

	if config.GitHubApp != nil {
		privateKey, err := loadPrivateKey(config.GitHubApp)
		if err != nil {
			return nil, err
		}

		return NewGitHubAppAuth(
			config.GitHubApp.AppID,
			config.GitHubApp.InstallationID,
			privateKey,
		)
	}

	return nil, errors.New("no authentication method provided")
}

// Helper to load private key from path or direct content
func loadPrivateKey(appConfig *config.GitHubApp) ([]byte, error) {
	if appConfig.PrivateKey != "" {
		return []byte(appConfig.PrivateKey), nil
	}

	if appConfig.PrivateKeyPath != "" {
		return os.ReadFile(appConfig.PrivateKeyPath)
	}

	return nil, errors.New("no private key provided")
}
