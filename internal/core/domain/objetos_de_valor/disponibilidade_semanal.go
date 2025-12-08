package objetosdevalor

type DisponibilidadeSemanal struct {
    PrestadorID  string
    DiaDaSemana  string // Ex: "SEGUNDA", "TERCA"
    Intervalos   []IntervaloDiario // Lista de horários de trabalho naquele dia
}

type IntervaloDiario struct {
    HoraInicio string // Ex: "09:00"
    HoraFim    string // Ex: "12:00"
}