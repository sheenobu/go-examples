package rest

import "time"

// A User is a user of our system
type User struct {
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	LastSeen time.Time `json:"last_seen"`
}
