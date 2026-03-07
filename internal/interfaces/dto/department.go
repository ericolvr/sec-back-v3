package dto

import "time"

type DepartmentRequest struct {
	CompanyID int64  `json:"company_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Active    bool   `json:"active" binding:"required"`
}

type DepartmentUpdateRequest struct {
	Name   string `json:"name" binding:"required"`
	Active bool   `json:"active" binding:"required"`
}

type DepartmentResponse struct {
	ID             int       `json:"id"`
	PartnerID      int       `json:"partner_id"`
	CompanyID      int       `json:"company_id"`
	Name           string    `json:"name"`
	Active         bool      `json:"active"`
	TotalEmployees int       `json:"total_employees"`
	CreatedAt      time.Time `json:"created_at"`
}
