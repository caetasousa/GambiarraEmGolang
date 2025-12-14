package request

import (
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorRequest struct {
	Nome        string   `json:"nome" binding:"required,min=3,max=100" example:"joao" swagger:"desc('Nome do prestador')"`
	Email       string   `json:"email" binding:"omitempty,email" example:"joao@email.com" swagger:"desc('Email do prestador')"`
	Telefone    string   `json:"telefone" binding:"required,min=8,max=15" example:"62999677481" swagger:"desc('Telefone do prestador')"`
	CatalogoIDs []string `json:"catalogo_ids" binding:"omitempty,dive,required" swagger:"desc('IDs dos serviços no catálogo oferecidos pelo prestador')"`
}

func (r *PrestadorRequest) ToPrestador(catalogos []domain.Catalogo) (*domain.Prestador, error) {
	return domain.NovoPrestador(r.Nome, r.Email, r.Telefone, catalogos)
}
