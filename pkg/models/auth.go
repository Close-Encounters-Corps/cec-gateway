package models

type AuthPhaseResult struct {

	// next url
	NextURL string `json:"next_url,omitempty"`

	// phase
	Phase int32 `json:"phase,omitempty"`

	// token
	Token string `json:"token,omitempty"`

	// user
	User *User `json:"user,omitempty"`
}
