package output

type BuscarPrestadorOutput struct {
	ID        string
	Nome      string
	Email     string
	Telefone  string
	Cpf       string
	Ativo     bool
	ImagemUrl string
	Catalogo  []CatalogoOutput
	Agenda    []AgendaDiariaOutput
}
type AgendaDiariaOutput struct {
	ID         string
	Data       string
	Intervalos []IntervaloDiarioOutput
}

type IntervaloDiarioOutput struct {
	ID         string
	HoraInicio string
	HoraFim    string
}
