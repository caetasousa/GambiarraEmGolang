package output

import (
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoOutput struct {
	ID             string
	Cliente        *domain.Cliente
	Prestador      *domain.Prestador
	Catalogo       *domain.Catalogo
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         domain.StatusDoAgendamento
	Notas          string
}
