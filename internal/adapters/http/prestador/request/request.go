package request

import (
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
)

type PrestadorRequest struct {
	Nome        string   `json:"nome" binding:"required"`
	Email       string   `json:"email" binding:"omitempty,email"`
	Telefone    string   `json:"telefone" binding:"required"`
	ServicosIDs []string `json:"servicos_ids" binding:"omitempty,dive,required"`
}

func (r *PrestadorRequest) ToPrestador() *domain.Prestador {
	return &domain.Prestador{
		ID:          xid.New().String(),
		Nome:        r.Nome,
		Email:       r.Email,
		Telefone:    r.Telefone,
		Ativo:       true,
		ServicosIDs: r.ServicosIDs,
	}
}
