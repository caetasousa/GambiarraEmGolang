package response_agendamento

import "meu-servico-agenda/internal/core/application/output"

type AgendamentoListPaginadoResponse struct {
	Data  []*AgendamentoResponse `json:"data"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
	Total int                    `json:"total"`
}

func ToAgendamentoListPaginadoResponse(outputs []*output.AgendamentoOutput, page, limit, total int) *AgendamentoListPaginadoResponse {
	responses := make([]*AgendamentoResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = NovoAgendamentoResponse(output)
	}

	return &AgendamentoListPaginadoResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}
}
