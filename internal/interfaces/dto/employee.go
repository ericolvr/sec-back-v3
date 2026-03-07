package dto

type EmployeeRequest struct {
	CompanyID    int64  `json:"company_id"`
	DepartmentID int64  `json:"department_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile"`
	Active       bool   `json:"active"`
}

type EmployeeResponse struct {
	ID             int    `json:"id"`
	CompanyID      int64  `json:"company_id"`
	DepartmentID   int64  `json:"department_id"`
	DepartmentName string `json:"department_name,omitempty"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Mobile         string `json:"mobile"`
	Active         bool   `json:"active"`
	CreatedAt      string `json:"created_at"`
}
