package dto

type EmployeeRequest struct {
	Name         string `json:"name"`
	DepartmentID int64  `json:"department_id"`
	Email        string `json:"email"`
	Active       bool   `json:"active"`
}

type EmployeeResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name,omitempty"`
	Email          string `json:"email"`
	Active         bool   `json:"active"`
	CreatedAt      string `json:"created_at"`
}
