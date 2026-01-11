package response

type DebugEchoResponse struct {
	Message string `json:"message"`
}

// DebugJWTResponse is the response for JWT generation endpoint
type DebugJWTResponse struct {
	Token  string `json:"token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserID string `json:"user_id" example:"user-12345"`
}

// DebugIDTokenResponse is the response for ID token generation endpoint
type DebugIDTokenResponse struct {
	IDToken   string `json:"id_token" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjFlOWdkazcifQ..."`
	ExpiresIn string `json:"expires_in" example:"3600"`
	UserID    string `json:"user_id" example:"user-12345"`
	TokenType string `json:"token_type" example:"Bearer"`
}
