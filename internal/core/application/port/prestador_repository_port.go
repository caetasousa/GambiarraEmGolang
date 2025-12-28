package port

import (
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorRepositorio interface {
	Salvar(prestador *domain.Prestador) error
	BuscarPorId(id string) (*domain.Prestador, error)
	BuscarPorCPF(cpf string) (*domain.Prestador, error)
	BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error)
	Atualizar(prestador *input.AlterarPrestadorInput) error
	Listar(input *input.PrestadorListInput) ([]*domain.Prestador, error)
	Contar(ativo bool) (int, error)
	AtualizarStatus(id string, ativo bool) error 
}
