package request_agendamento

type AgendamentoRequest struct {
	ClienteID      string `json:"cliente_id" binding:"required" swagger:"desc('ID do cliente que está solicitando o agendamento')"`
	PrestadorID    string `json:"prestador_id" binding:"required" swagger:"desc('ID do prestador que irá atender')"`
	CatalogoID     string `json:"catalogo_id" binding:"required" swagger:"desc('ID do serviço do catálogo que será agendado')"`
	DataHoraInicio string `json:"data_hora_inicio" binding:"required,datetime=2006-01-02T15:04:05Z07:00" example:"2025-01-03T08:00:00Z"`
	Notas          string `json:"notas,omitempty" binding:"omitempty,max=500" swagger:"desc('Notas ou observações do cliente sobre o agendamento')"`
}

// func (ag *AgendamentoRequest) ToAgendamento() (*domain.Agendamento, error) {
// 	// Convertendo DataHoraInicio de string para time.Time
// 	dataHoraInicio, err := time.Parse(time.RFC3339, ag.DataHoraInicio)
// 	if err != nil {
// 		return nil, errors.New("formato de data/hora inválido")
// 	}

// 	// Usando o método NovoAgendamento da camada de domínio para criar o agendamento
// 	agendamento, err := domain.NovoAgendamento(
// 		ag.ClienteID,
// 		ag.PrestadorID,
// 		ag.CatalogoID,
// 		dataHoraInicio,
// 		dataHoraInicio.Add(time.Duration(30)*time.Minute),
// 		ag.Notas,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return agendamento, nil
// }
