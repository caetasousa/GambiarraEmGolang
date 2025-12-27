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
	// Converte map para slice
	todos := make([]*domain.Prestador, 0, len(r.storage))
	for _, p := range r.storage {
		todos = append(todos, p)
	}

	// Ordena por ID (simulando ORDER BY created_at DESC do banco real)
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID > todos[j].ID
	})

	// Calcula offset
	offset := (input.Page - 1) * input.Limit

	// Se offset estiver além do tamanho, retorna vazio
	if offset >= len(todos) {
		return []*domain.Prestador{}, nil
	}

	// Calcula o fim da página
	fim := offset + input.Limit
	if fim > len(todos) {
		fim = len(todos)
	}

	// Retorna slice paginado
	return todos[offset:fim], nil
}

func (r *FakePrestadorRepositorio) Contar() (int, error) {
	return len(r.storage), nil
}