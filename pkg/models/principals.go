package models

import "time"

type User struct {

	// id
	ID int64 `json:"id,omitempty"`

	// principal
	Principal *Principal `json:"principal,omitempty"`
}

type Principal struct {

	// admin
	Admin bool `json:"admin"`

	// created on
	CreatedOn *time.Time `json:"created_on"`

	// id
	ID int64 `json:"id"`

	// last login
	LastLogin *time.Time `json:"last_login"`

	// state
	State string `json:"state"`
}
