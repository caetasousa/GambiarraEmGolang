package response_agendamento

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type PrestadorInfo struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
}

type ServicoInfo struct {
	Nome    string `json:"nome"`
	Duracao int    `json:"duracao"`
	Preco   int    `json:"preco"`
}

type AgendamentoResponse struct {
	ID         string                     `json:"id"`
	Prestador  PrestadorInfo              `json:"prestador"`
	Servico    ServicoInfo                `json:"servico"`
	DataInicio time.Time                  `json:"data_inicio"`
	DataFim    time.Time                  `json:"data_fim"`
	Status     domain.StatusDoAgendamento `json:"status"`
	Notas      string                     `json:"notas,omitempty"`
}

func NovoAgendamentoResponse(a *output.AgendamentoOutput) *AgendamentoResponse {
	return &AgendamentoResponse{
		ID: a.ID,
		Prestador: PrestadorInfo{
			Nome:     a.PrestadorNome,
			Telefone: a.PrestadorTel,
		},
		Servico: ServicoInfo{
			Nome:    a.ServicoNome,
			Duracao: a.ServicoDuracao,
			Preco:   a.ServicoPreco,
		},
		DataInicio: a.DataHoraInicio,
		DataFim:    a.DataHoraFim,
		Status:     a.Status,
		Notas:      a.Notas,
	}
}
