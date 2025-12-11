package request

import (
	"fmt"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
)

type PrestadorRequest struct {
	Nome        string   `json:"nome" binding:"required,min=3,max=100" swagger:"desc('Nome do prestador')"`
	Email       string   `json:"email" binding:"omitempty,email" swagger:"desc('Email do prestador')"`
	Telefone    string   `json:"telefone" binding:"required,min=8,max=15" swagger:"desc('Telefone do prestador')"`
	CatalogoIDs []string `json:"catalogo_ids" binding:"omitempty,dive,required" swagger:"desc('IDs dos serviços no catálogo oferecidos pelo prestador')"`
}

func (r *PrestadorRequest) ToPrestador(catalogoRepo port.CatalogoRepositorio) (*domain.Prestador, error) {
	catalogos, err := r.ValidateCatalogoIDs(catalogoRepo)
	if err != nil {
		return nil, err
	}

	return &domain.Prestador{
		ID:          xid.New().String(),
		Nome:        r.Nome,
		Email:       r.Email,
		Telefone:    r.Telefone,
		Ativo:       true,
		Catalogo: catalogos,
	}, nil
}

// ValidateCatalogoIDs verifica se todos os IDs existem no repositório
func (r *PrestadorRequest) ValidateCatalogoIDs(catalogoRepo port.CatalogoRepositorio) ([]domain.Catalogo, error) {
	catalogos := []domain.Catalogo{}
	for _, id := range r.CatalogoIDs {
		cat, err := catalogoRepo.BuscarPorId(id)
		if err != nil {
			return nil, fmt.Errorf("erro ao consultar catálogo: %w", err)
		}
		if cat == nil {
			return nil, fmt.Errorf("o catálogo com ID '%s' não existe", id)
		}
		catalogos = append(catalogos, *cat)
	}
	return catalogos, nil
}
