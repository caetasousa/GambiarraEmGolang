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
	agendaDiariaRepo port.AgendaDiariaRepositorio
}

func NovaPrestadorService(pr port.PrestadorRepositorio, cr port.CatalogoRepositorio, ad port.AgendaDiariaRepositorio) *PrestadorService {
	return &PrestadorService{
		prestadorRepo: pr,
		catalogoRepo:  cr,
		agendaDiariaRepo: ad,
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

func (s *PrestadorService) AdicionarAgenda(prestadorID string, req *request.AgendaDiariaRequest) error {

	// 1️⃣ Buscar agregado
	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		return err
	}
	if prestador == nil {
		return fmt.Errorf("prestador não encontrado")
	}

	// 2️⃣ Criar entidade de domínio
	agenda, err := req.ToAgendaDiaria()
	if err != nil {
		return err
	}

	// 3️⃣ Delegar decisão ao domínio
	if err := prestador.AdicionarAgenda(agenda); err != nil {
		return err
	}

	// 4️⃣ Persistir agenda (tabela separada)
	if err := s.agendaDiariaRepo.Salvar(agenda); err != nil {
		return err
	}

	// 5️⃣ Persistir agregado inteiro
	return s.prestadorRepo.Salvar(prestador)
}
