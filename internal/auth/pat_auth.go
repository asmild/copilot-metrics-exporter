package auth

// PATAuth implements AuthProvider for Personal Access Tokens
type PATAuth struct {
	Token string
}

// NewPATAuth creates a new PAT auth provider
func NewPATAuth(token string) *PATAuth {
	return &PATAuth{Token: token}
}

// Token returns the PAT
func (a *PATAuth) Token() (string, error) {
	return a.Token, nil
}
