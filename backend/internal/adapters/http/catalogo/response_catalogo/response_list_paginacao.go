package response_catalogo

import "meu-servico-agenda/internal/core/application/output"

type CatalogoListResponse struct {
	Data  []*output.CatalogoOutput `json:"data"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
	Total int                      `json:"total"`
}
