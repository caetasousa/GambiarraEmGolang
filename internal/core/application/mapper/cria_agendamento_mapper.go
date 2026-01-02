package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

func NovoAgendamentoOutput(a *domain.Agendamento) *output.AgendamentoOutput {
	return &output.AgendamentoOutput{
		ID:             a.ID,
		Cliente:        a.Cliente,
		Prestador:      a.Prestador,
		Catalogo:       a.Catalogo,
		DataHoraInicio: a.DataHoraInicio,
		DataHoraFim:    a.DataHoraFim,
		Status:         a.Status,
		Notas:          a.Notas,
	}
}

func BuscaAgendamentoClienteData(agendamentos []*domain.Agendamento) []*output.AgendamentoOutput {
	outputs := make([]*output.AgendamentoOutput, len(agendamentos))
	for i, ag := range agendamentos {
		outputs[i] = NovoAgendamentoOutput(ag)
	}
	return outputs
}
