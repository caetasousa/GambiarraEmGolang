package domain

type PrestadorDeServico struct {
	ID            string
	Nome          string
	Email         string
	Telefone      string
	Ativo         bool
	ServicosIDs []string
}
