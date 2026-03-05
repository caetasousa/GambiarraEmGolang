package response_prestador

import (
	"meu-servico-agenda/internal/core/application/output"
)


type AgendaDiariaResponse struct {
	ID         string                    `json:"id"`
	Data       string                    `json:"data"`
	Intervalos []IntervaloDiarioResponse `json:"intervalos"`
}

type IntervaloDiarioResponse struct {
	ID         string `json:"id"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

func FromAgendaDiariaOutput(o output.AgendaDiariaOutput) AgendaDiariaResponse {
	intervalos := make([]IntervaloDiarioResponse, len(o.Intervalos))
	for i, it := range o.Intervalos {
		intervalos[i] = IntervaloDiarioResponse{
			ID:         it.ID,
			HoraInicio: it.HoraInicio,
			HoraFim:    it.HoraFim,
		}
	}

	return AgendaDiariaResponse{
		ID:         o.ID,
		Data:       o.Data,
		Intervalos: intervalos,
	}
}
