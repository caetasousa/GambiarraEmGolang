package request_catalogo

import "meu-servico-agenda/internal/core/application/input"

type CatalogoUpdateRequest struct {
	Nome          string `json:"nome" binding:"required,min=3,max=100" example:"Limpesa de Peele" swagger:"desc('Nome do serviço')"`
	DuracaoPadrao int    `json:"duracao_padrao" binding:"required" example:"40" swagger:"desc('Duração padrão do serviço em minutos')"`
	Preco         int    `json:"preco" binding:"required" example:"10000" swagger:"desc('Preço do serviço em centavos')"`
	ImagemUrl     string `json:"image_url" binding:"required,url" example:"https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/bb515383d2f6ef76.jpg"`
	Categoria     string `json:"categoria" binding:"required,min=3,max=50" example:"Estética Facial" swagger:"desc('Categoria do serviço')"`
}

func (cr *CatalogoUpdateRequest) ToCatalogoUpdateInput() *input.CatalogoUpdateInput {
	return &input.CatalogoUpdateInput{
		Nome:          cr.Nome,
		DuracaoPadrao: cr.DuracaoPadrao,
		Preco:         cr.Preco,
		Categoria:     cr.Categoria,
		ImagemUrl:     cr.ImagemUrl,
	}
}
