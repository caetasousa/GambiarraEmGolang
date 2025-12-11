package domain

// AgendaDiaria: A estrutura central que define a agenda para uma data específica.
type AgendaDiaria struct {
	Id         string
	Prestador  Prestador
	Data       string            // Ex: "2026-01-15" (A data em que o trabalho ocorre)
	Intervalos []IntervaloDiario // Lista de horários disponíveis naquele dia
}

// IntervaloDiario: O bloco de tempo (Ex: 09:00 - 12:00)
type IntervaloDiario struct {
	Id         string
	HoraInicio string // Ex: "09:00"
	HoraFim    string // Ex: "12:00"
}
