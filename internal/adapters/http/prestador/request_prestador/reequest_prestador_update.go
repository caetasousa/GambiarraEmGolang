package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type PrestadorUpdateRequest struct {
	Nome        string   `json:"nome" binding:"required,min=3,max=100" example:"joao" swagger:"desc('Nome do prestador')"`
	Email       string   `json:"email" binding:"omitempty,email" example:"joao@email.com" swagger:"desc('Email do prestador')"`
	Telefone    string   `json:"telefone" binding:"required,min=8,max=15" example:"62999677481" swagger:"desc('Telefone do prestador')"`
	ImagemUrl   string   `json:"image_url" binding:"required,url" example:"https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/bb515383d2f6ef76.jpg"`
	CatalogoIDs []string `json:"catalogo_ids" binding:"omitempty,dive,required" swagger:"desc('IDs dos serviços no catálogo oferecidos pelo prestador')"`
}

func (r *PrestadorUpdateRequest) ToAlterarPrestadorInput() (*input.AlterarPrestadorInput) {
	return &input.AlterarPrestadorInput{
		Nome:        r.Nome,
		Email:       r.Email,
		Telefone:    r.Telefone,
		ImagemUrl:   r.ImagemUrl,
		CatalogoIDs: r.CatalogoIDs,
	}
}
