package domain

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	//Valida Prestador
	ErrAgendaDuplicada          = errors.New("agenda duplicada")
	ErrPrestadorInativo         = errors.New("prestador inativo")
	ErrPrestadorDeveTerCatalogo = errors.New("prestador deve ter ao menos um catálogo de serviços")
	ErrAgendaNaoEncontrada = errors.New("agenda não encontrada para esta data")

	//Valida Catalogo
	ErrDuracaoInvalida   = errors.New("duração padrão inválida")
	ErrPrecoInvalido     = errors.New("preço inválido")
	ErrNomeInvalido      = errors.New("Nome invalido")
	ErrCategoriaInvalida = errors.New("Categoria invalida")

	//Validaa Agendamento
	ErrHoraInicialMenorQueFinal = errors.New("horário início deve ser antes do fim")

	//Valida Agenda Diaria
	ErrAgendaSemIntervalos      = errors.New("agenda deve conter ao menos um intervalo")
	ErrIntervaloHorarioInvalido = errors.New("hora início deve ser menor que hora fim")
	ErrDataEstaNoPassado        = errors.New("Esta data está passado")
	ErrIntervalosSesobrepoe     = errors.New("intervalos de horário não podem se sobrepor")
)
