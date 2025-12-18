package request_catalogo

import "meu-servico-agenda/internal/core/application/command"

type CatalogoRequest struct {
	Nome          string  `json:"nome" binding:"required,min=3,max=100" example:"Tecnico de Redes" swagger:"desc('Nome do serviço')"`
	DuracaoPadrao int     `json:"duracao_padrao" binding:"required" example:"20" swagger:"desc('Duração padrão do serviço em minutos')"`
	Preco         float64 `json:"preco" binding:"required" example:"10000" swagger:"desc('Preço do serviço em centavos')"`
	Categoria     string  `json:"categoria" binding:"required,min=3,max=50" example:"Redes" swagger:"desc('Categoria do serviço')"`
}

func (cr *CatalogoRequest) ToCommand() *command.CatalogoCommand {
	return &command.CatalogoCommand{
		Nome:          cr.Nome,
		DuracaoPadrao: cr.DuracaoPadrao,
		Preco:         cr.Preco,
		Categoria:     cr.Categoria,
	}
}
