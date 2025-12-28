package repository

import (
	"database/sql"
	"fmt"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type FakeAgendaDiariaRepositorio struct {
	storage map[string]*domain.AgendaDiaria // chave: "prestadorID:data" → agenda
}

func NovoFakeAgendaDiariaRepositorio() port.AgendaDiariaRepositorio {
	return &FakeAgendaDiariaRepositorio{
		storage: make(map[string]*domain.AgendaDiaria),
	}
}

func (r *FakeAgendaDiariaRepositorio) Salvar(agenda *domain.AgendaDiaria, prestadorId string) error {
	// Cria chave única: prestadorID + data
	chave := fmt.Sprintf("%s:%s", prestadorId, agenda.Data)
	r.storage[chave] = agenda
	return nil
}

func (r *FakeAgendaDiariaRepositorio) BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error) {
	// Cria chave única: prestadorID + data
	chave := fmt.Sprintf("%s:%s", prestadorID, data)
	
	// Busca agenda
	agenda, exists := r.storage[chave]
	if !exists {
		return nil, sql.ErrNoRows
	}

	return agenda, nil
}

func (r *FakeAgendaDiariaRepositorio) AtualizarAgenda(agenda *domain.AgendaDiaria, prestadorID string) error {
	// Cria chave única: prestadorID + data
	chave := fmt.Sprintf("%s:%s", prestadorID, agenda.Data)

	// Verifica se agenda existe
	if _, exists := r.storage[chave]; !exists {
		return sql.ErrNoRows
	}

	// Atualiza (substitui completamente)
	r.storage[chave] = agenda

	return nil
}