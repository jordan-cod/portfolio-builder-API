package models

type RegisterInput struct {
	Name     string `json:"name" example:"João"`
	Email    string `json:"email" example:"joao@email.com"`
	Password string `json:"password" example:"minhasenha123"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"Usuário cadastrado com sucesso"`
	APIKey  string `json:"api_key" example:"meu-token-api-gerado"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Erro ao validar os dados"`
}
