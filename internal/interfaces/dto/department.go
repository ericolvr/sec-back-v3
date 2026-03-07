package dto

import "time"

type DepartmentRequest struct {
	Name string `json:"name"`
}

type DepartmentResponse struct {
	ID             int       `json:"id"`
	PartnerID       int       `json:"tenant_id"`
	Name           string    `json:"name"`
	TotalEmployees int       `json:"total_employees"`
	CreatedAt      time.Time `json:"created_at"`
}
