package response_prestador

import "meu-servico-agenda/internal/core/application/output"

type PrestadorListResponse struct {
	Data  []*output.BuscarPrestadorOutput `json:"data"`
	Page  int                             `json:"page"`
	Limit int                             `json:"limit"`
	Total int                             `json:"total"`
}
