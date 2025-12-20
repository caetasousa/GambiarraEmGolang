package mapper

import (
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoOutput struct {
	ID             string
	PrestadorNome  string
	PrestadorTel   string
	ServicoNome    string
	ServicoDuracao int
	ServicoPreco   int
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         domain.StatusDoAgendamento  
	Notas          string
}

func NovoAgendamentoOutput(a *domain.Agendamento) *AgendamentoOutput {
	return &AgendamentoOutput{
		ID:             a.ID,
		PrestadorNome:  a.Prestador.Nome,
		PrestadorTel:   a.Prestador.Telefone,
		ServicoNome:    a.Catalogo.Nome,
		ServicoDuracao: a.Catalogo.DuracaoPadrao,
		ServicoPreco:   a.Catalogo.Preco,
		DataHoraInicio: a.DataHoraInicio,
		DataHoraFim:    a.DataHoraFim,
		Status:         a.Status, 
		Notas:          a.Notas,
	}
}

