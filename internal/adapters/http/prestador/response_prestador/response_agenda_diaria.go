package response_prestador

import (
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendaDiariaResponse struct {
	Id         string                    `json:"id"`
	Data       string                    `json:"data"`
	DiaSemana  string                    `json:"dia_semana"`
	Intervalos []IntervaloDiarioResponse `json:"intervalos"`
}

type IntervaloDiarioResponse struct {
	Id         string `json:"id"`
	HoraInicio string `json:"hora_inicio"`
	HoraFim    string `json:"hora_fim"`
}

var diasSemana = map[time.Weekday]string{
	time.Monday:    "segunda",
	time.Tuesday:   "terça",
	time.Wednesday: "quarta",
	time.Thursday:  "quinta",
	time.Friday:    "sexta",
	time.Saturday:  "sábado",
	time.Sunday:    "domingo",
}

func FromAgendaDiaria(ag *domain.AgendaDiaria) AgendaDiariaResponse {

	// Parse da data vinda do domínio
	data, _ := time.Parse("2006-01-02", ag.Data)

	// Monta intervalos
	intervalos := make([]IntervaloDiarioResponse, len(ag.Intervalos))
	for i, it := range ag.Intervalos {
		intervalos[i] = IntervaloDiarioResponse{
			Id:         it.Id,
			HoraInicio: it.HoraInicio.Format("15:04"),
			HoraFim:    it.HoraFim.Format("15:04"),
		}
	}

	return AgendaDiariaResponse{
		Id:         ag.Id,
		Data:       ag.Data,
		DiaSemana:  diasSemana[data.Weekday()],
		Intervalos: intervalos,
	}
}
