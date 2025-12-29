package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"sort"

	"github.com/klassmann/cpfcnpj"
)

type FakePrestadorRepositorio struct {
	storage      map[string]*domain.Prestador
	catalogoRepo port.CatalogoRepositorio
}

func NovoFakePrestadorRepositorio(catalogoRepo port.CatalogoRepositorio) port.PrestadorRepositorio {
	return &FakePrestadorRepositorio{
		storage:      make(map[string]*domain.Prestador),
		catalogoRepo: catalogoRepo,
	}
}

func (r *FakePrestadorRepositorio) Salvar(prestador *domain.Prestador) error {
	r.storage[prestador.ID] = prestador
	return nil
}

func (r *FakePrestadorRepositorio) BuscarPorId(id string) (*domain.Prestador, error) {
	prestador := r.storage[id]
	if prestador == nil {
		return nil, errors.New("não encontrado")
	}
	return prestador, nil
}

func (r *FakePrestadorRepositorio) BuscarPorCPF(cpf string) (*domain.Prestador, error) {
	cpf = cpfcnpj.Clean(cpf)
	for _, p := range r.storage {
		if cpfcnpj.Clean(p.Cpf) == cpf {
			return p, nil
		}
	}
	return nil, nil
}

func (r *FakePrestadorRepositorio) BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error) {
	prestador, ok := r.storage[prestadorID]
	if !ok {
		return nil, nil
	}

	for _, agenda := range prestador.Agenda {
		if agenda.Data == data {
			return &agenda, nil
		}
	}

	return nil, nil
}

func (r *FakePrestadorRepositorio) Atualizar(input *input.AlterarPrestadorInput) error {
	// 1️⃣ Verifica se o prestador existe
	prestador, exists := r.storage[input.Id]
	if !exists {
		return sql.ErrNoRows
	}

	// 2️⃣ Valida se os catálogos existem
	for _, catalogoID := range input.CatalogoIDs {
		_, err := r.catalogoRepo.BuscarPorId(catalogoID)
		if err != nil {
			return fmt.Errorf("catálogo %s não existe", catalogoID)
		}
	}

	// 3️⃣ Atualiza os campos editáveis
	prestador.Nome = input.Nome
	prestador.Email = input.Email
	prestador.Telefone = input.Telefone
	prestador.ImagemUrl = input.ImagemUrl

	// 4️⃣ Atualiza os catálogos
	novos := make([]domain.Catalogo, len(input.CatalogoIDs))
	for i, catalogoID := range input.CatalogoIDs {
		catalogo, _ := r.catalogoRepo.BuscarPorId(catalogoID)
		novos[i] = *catalogo
	}
	prestador.Catalogo = novos

	// 5️⃣ Salva de volta
	r.storage[input.Id] = prestador

	return nil
}

func (r *FakePrestadorRepositorio) Listar(input *input.PrestadorListInput) ([]*domain.Prestador, error) {
	// ✅ Sempre filtra por status (obrigatório)
	todos := make([]*domain.Prestador, 0, len(r.storage))
	for _, p := range r.storage {
		if p.Ativo == input.Ativo {
			todos = append(todos, p)
		}
	}

	// Ordena por ID (simulando ORDER BY created_at DESC)
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID > todos[j].ID
	})

	// Calcula offset
	offset := (input.Page - 1) * input.Limit

	if offset >= len(todos) {
		return []*domain.Prestador{}, nil
	}

	fim := offset + input.Limit
	if fim > len(todos) {
		fim = len(todos)
	}

	return todos[offset:fim], nil
}

func (r *FakePrestadorRepositorio) Contar(ativo bool) (int, error) {
	// ✅ Conta apenas os que correspondem ao filtro
	count := 0
	for _, p := range r.storage {
		if p.Ativo == ativo {
			count++
		}
	}
	
	return count, nil
}

func (r *FakePrestadorRepositorio) AtualizarStatus(id string, ativo bool) error {
	prestador, exists := r.storage[id]
	if !exists {
		return sql.ErrNoRows
	}

	prestador.Ativo = ativo
	r.storage[id] = prestador

	return nil
}

func (r *FakeAgendaDiariaRepositorio) DeletarAgenda(prestadorID string, data string) error {
	chave := fmt.Sprintf("%s:%s", prestadorID, data)
	
	if _, exists := r.storage[chave]; !exists {
		return sql.ErrNoRows
	}

	delete(r.storage, chave)
	return nil
}

func (r *FakePrestadorRepositorio) BuscarPrestadoresDisponiveisPorData(data string, page, limit int) ([]*domain.Prestador, error) {
	// Filtra prestadores ativos que têm agenda na data
	disponiveis := make([]*domain.Prestador, 0)
	
	for _, p := range r.storage {
		// Só considera prestadores ativos
		if !p.Ativo {
			continue
		}
		
		// Verifica se tem agenda na data
		temAgenda := false
		for _, agenda := range p.Agenda {
			if agenda.Data == data {
				temAgenda = true
				break
			}
		}
		
		if temAgenda {
			disponiveis = append(disponiveis, p)
		}
	}

	// Ordena por ID (simulando ORDER BY created_at DESC)
	sort.Slice(disponiveis, func(i, j int) bool {
		return disponiveis[i].ID > disponiveis[j].ID
	})

	// Calcula offset
	offset := (page - 1) * limit

	if offset >= len(disponiveis) {
		return []*domain.Prestador{}, nil
	}

	fim := offset + limit
	if fim > len(disponiveis) {
		fim = len(disponiveis)
	}

	return disponiveis[offset:fim], nil
}

// ContarPrestadoresDisponiveisPorData conta quantos prestadores ativos têm agenda na data informada
func (r *FakePrestadorRepositorio) ContarPrestadoresDisponiveisPorData(data string) (int, error) {
	count := 0
	
	for _, p := range r.storage {
		// Só conta prestadores ativos
		if !p.Ativo {
			continue
		}
		
		// Verifica se tem agenda na data
		for _, agenda := range p.Agenda {
			if agenda.Data == data {
				count++
				break // Já encontrou, não precisa continuar
			}
		}
	}
	
	return count, nil
}