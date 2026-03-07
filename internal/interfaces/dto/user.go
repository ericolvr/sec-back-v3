package dto

import "time"

type UserRequest struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Type   int    `json:"type"` // 1 = master, 2 = client
	Active bool   `json:"active"`
}

type UserPasswordResetRequest struct {
	Mobile string `json:"mobile"`
}

type LoginRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	PartnerID  int       `json:"partner_id"`
	Name      string    `json:"name"`
	Mobile    string    `json:"mobile"`
	Type      int       `json:"type"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}
