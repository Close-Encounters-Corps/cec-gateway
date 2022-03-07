package models

type User struct {

	// id
	ID int64 `json:"id,omitempty"`

	// principal
	Principal *Principal `json:"principal,omitempty"`
}

type Principal struct {

	// admin
	Admin bool `json:"admin,omitempty"`

	// created on
	CreatedOn string `json:"created_on,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// last login
	LastLogin string `json:"last_login,omitempty"`

	// state
	State string `json:"state,omitempty"`
}
