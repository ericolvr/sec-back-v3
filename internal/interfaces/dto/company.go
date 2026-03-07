package dto

type CompanyRequest struct {
	Name   string `json:"name"`
	CNPJ   string `json:"cnpj"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type CompanyResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CNPJ      string `json:"cnpj"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
