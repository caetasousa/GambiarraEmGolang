package domain

import objetosdevalor "meu-servico-agenda/internal/core/domain/objetos_de_valor"


type PrestadorDeServico struct {
	ID              string
	Nome            string
	Email           string
	Telefone        string
	Disponibilidade objetosdevalor.DisponibilidadeSemanal
}
