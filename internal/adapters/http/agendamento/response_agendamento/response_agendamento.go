package response_agendamento

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type PrestadorInfo struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	CPF      string `json:"cpf"`
	Email    string `json:"email,omitempty"`
	Telefone string `json:"telefone"`
	Ativo    bool   `json:"ativo"`
}

type ClienteInfo struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type ServicoInfo struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Duracao   int    `json:"duracao"`
	Preco     int    `json:"preco"`
	Categoria string `json:"categoria"`
}

type AgendamentoResponse struct {
	ID         string                     `json:"id"`
	Cliente    ClienteInfo                `json:"cliente"`
	Prestador  PrestadorInfo              `json:"prestador"`
	Servico    ServicoInfo                `json:"servico"`
	DataInicio time.Time                  `json:"data_inicio"`
	DataFim    time.Time                  `json:"data_fim"`
	Status     domain.StatusDoAgendamento `json:"status"`
	Notas      string                     `json:"notas,omitempty"`
}

func NovoAgendamentoResponse(a *output.AgendamentoOutput) *AgendamentoResponse {
	return &AgendamentoResponse{
		ID: a.ID,
		Cliente: ClienteInfo{
			ID:       a.Cliente.ID,
			Nome:     a.Cliente.Nome,
			Email:    a.Cliente.Email,
			Telefone: a.Cliente.Telefone,
		},
		Prestador: PrestadorInfo{
			ID:       a.Prestador.ID,
			Nome:     a.Prestador.Nome,
			CPF:      a.Prestador.Cpf,
			Email:    a.Prestador.Email,
			Telefone: a.Prestador.Telefone,
			Ativo:    a.Prestador.Ativo,
		},
		Servico: ServicoInfo{
			ID:        a.Catalogo.ID,
			Nome:      a.Catalogo.Nome,
			Duracao:   a.Catalogo.DuracaoPadrao,
			Preco:     a.Catalogo.Preco,
			Categoria: a.Catalogo.Categoria,
		},
		DataInicio: a.DataHoraInicio,
		DataFim:    a.DataHoraFim,
		Status:     a.Status,
		Notas:      a.Notas,
	}
}