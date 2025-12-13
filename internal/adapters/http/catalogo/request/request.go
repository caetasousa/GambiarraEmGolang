package request

import (
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoRequest struct {
	Nome          string  `json:"nome" binding:"required,min=3,max=100" swagger:"desc('Nome do serviço')"`
	DuracaoPadrao int     `json:"duracao_padrao" binding:"required" swagger:"desc('Duração padrão do serviço em minutos')"`
	Preco         float64 `json:"preco" binding:"required" swagger:"desc('Preço do serviço em centavos')"`
	Categoria     string  `json:"categoria" binding:"required,min=3,max=50" swagger:"desc('Categoria do serviço')"`
}

func (cr *CatalogoRequest) ToCatalogo() (*domain.Catalogo, error) {
	return domain.NovoCatalogo(
		cr.Nome,
		cr.DuracaoPadrao,
		cr.Preco,
		cr.Categoria,
	)
}
