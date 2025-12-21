package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)
func NovoAgendamentoOutput(a *domain.Agendamento) *output.AgendamentoOutput{
	return &output.AgendamentoOutput{
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
