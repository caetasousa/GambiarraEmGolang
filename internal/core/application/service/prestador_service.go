package service

import (
	"fmt"
	"meu-servico-agenda/internal/adapters/http/prestador/request"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorService struct {
	prestadorRepo port.PrestadorRepositorio
	catalogoRepo  port.CatalogoRepositorio
}

func NovaPrestadorService(pr port.PrestadorRepositorio, cr port.CatalogoRepositorio) *PrestadorService {
	return &PrestadorService{
		prestadorRepo: pr,
		catalogoRepo:  cr,
	}
}

// Cadastra cria um prestador validando existência dos catálogos
// Cadastra cria um prestador completo a partir do DTO
func (s *PrestadorService) Cadastra(req *request.PrestadorRequest) (*domain.Prestador, error) {
	// 1️⃣ Validar/consultar catálogos existentes
	catalogos := []domain.Catalogo{}
	for _, id := range req.CatalogoIDs {
		c, err := s.catalogoRepo.BuscarPorId(id)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, fmt.Errorf("catálogo '%s' não existe", id)
		}
		catalogos = append(catalogos, *c)
	}

	// 2️⃣ Criar domínio com regras de negócio (ID gerado, catálogo obrigatório)
	prestador, err := req.ToPrestador(catalogos)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Persistir domínio
	if err := s.prestadorRepo.Salvar(prestador); err != nil {
		return nil, err
	}

	return prestador, nil
}
