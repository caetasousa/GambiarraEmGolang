package output

type CriarPrestadorOutput struct {
	ID        string           `json:"id"`
	Nome      string           `json:"nome"`
	Email     string           `json:"email"`
	Telefone  string           `json:"telefone"`
	Cpf       string           `json:"cpf"`
	ImagemUrl string           `json:"imagem_url"`
	Ativo     bool             `json:"ativo"`
	Catalogo  []CatalogoOutput `json:"catalogo"`
}
