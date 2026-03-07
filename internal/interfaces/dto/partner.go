package dto

type PartnerRequest struct {
	Name   string `json:"name" binding:"required"`
	CNPJ   string `json:"cnpj" binding:"required"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

type PartnerResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CNPJ      string `json:"cnpj"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}
