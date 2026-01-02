package response_agendamento

import "meu-servico-agenda/internal/core/application/output"

type BuscaClienteDataResponse struct {
	Data []*AgendamentoResponse `json:"data"`
}

func ToBuscaClienteDataResponse(outputs []*output.AgendamentoOutput) *BuscaClienteDataResponse {
	responses := make([]*AgendamentoResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = NovoAgendamentoResponse(output)
	}

	return &BuscaClienteDataResponse{
		Data: responses,
	}
}
