package response_agendamento

import "meu-servico-agenda/internal/core/application/output"

type BuscaDataResponse struct {
	Data []*AgendamentoResponse `json:"data"`
}

func ToBuscaDataResponse(outputs []*output.AgendamentoOutput) *BuscaDataResponse {
	responses := make([]*AgendamentoResponse, len(outputs))
	for i, output := range outputs {
		responses[i] = NovoAgendamentoResponse(output)
	}

	return &BuscaDataResponse{
		Data: responses,
	}
}
