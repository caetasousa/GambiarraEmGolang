package request

import (
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoRequest struct {
	Nome          string  `json:"nome" binding:"required,min=3,max=100" example:"Tecnico de Redes" swagger:"desc('Nome do serviço')"`
	DuracaoPadrao int     `json:"duracao_padrao" binding:"required" example:"20" swagger:"desc('Duração padrão do serviço em minutos')"`
	Preco         float64 `json:"preco" binding:"required" example:"10000" swagger:"desc('Preço do serviço em centavos')"`
	Categoria     string  `json:"categoria" binding:"required,min=3,max=50" example:"Redes" swagger:"desc('Categoria do serviço')"`
}

func (cr *CatalogoRequest) ToCatalogo() (*domain.Catalogo, error) {
	return domain.NovoCatalogo(
		cr.Nome,
		cr.DuracaoPadrao,
		cr.Preco,
		cr.Categoria,
	)
}
