package request_agendamento

import (
	"errors"
	"meu-servico-agenda/internal/core/application/input"
	"time"
)

type AgendamentoRequest struct {
	ClienteID      string `json:"cliente_id" binding:"required" swagger:"desc('ID do cliente que está solicitando o agendamento')"`
	PrestadorID    string `json:"prestador_id" binding:"required" swagger:"desc('ID do prestador que irá atender')"`
	CatalogoID     string `json:"catalogo_id" binding:"required" swagger:"desc('ID do serviço do catálogo que será agendado')"`
	DataHoraInicio string `json:"data_hora_inicio" binding:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2025-01-03T08:00:00Z"`
	Notas          string `json:"notas,omitempty" binding:"omitempty,max=500" swagger:"desc('Notas ou observações do cliente sobre o agendamento')"`
}

func (ag *AgendamentoRequest) ToAgendamento() (*input.CadastrarAgendamentoInput, error) {
	// Convertendo DataHoraInicio de string para time.Time
	dataHoraInicio, err := time.Parse(time.RFC3339, ag.DataHoraInicio)
	if err != nil {
		return nil, errors.New("formato de data/hora inválido")
	}

	agendamento := input.CadastrarAgendamentoInput{
		ClienteID:      ag.ClienteID,
		PrestadorID:    ag.PrestadorID,
		CatalogoID:     ag.CatalogoID,
		DataHoraInicio: dataHoraInicio,
		Notas:          ag.Notas,
	}

	return &agendamento, nil
}
