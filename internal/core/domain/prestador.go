package domain

type Prestador struct {
	ID            string
	Nome          string
	Email         string
	Telefone      string
	Ativo         bool
	ServicosIDs []string
}
