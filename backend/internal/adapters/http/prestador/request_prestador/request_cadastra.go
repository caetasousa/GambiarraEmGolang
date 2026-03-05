package request_prestador

import (
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type PrestadorRequest struct {
	Nome        string   `json:"nome" binding:"required,min=3,max=100" example:"joao" swagger:"desc('Nome do prestador')"`
	Cpf         string   `json:"cpf" binding:"required,len=11,numeric" example:"04423258196"`
	Email       string   `json:"email" binding:"omitempty,email" example:"joao@email.com" swagger:"desc('Email do prestador')"`
	Telefone    string   `json:"telefone" binding:"required,min=8,max=15" example:"62999677481" swagger:"desc('Telefone do prestador')"`
	Senha       string   `json:"senha" binding:"required,min=8" example:"senha123" swagger:"desc('Senha do prestador (mínimo 8 caracteres)')"`
	ImagemUrl   string   `json:"image_url" binding:"required,url" example:"https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/bb515383d2f6ef76.jpg"`
	CatalogoIDs []string `json:"catalogo_ids" binding:"required,min=1" swagger:"desc('IDs dos serviços no catálogo oferecidos pelo prestador')"`
}

func (r *PrestadorRequest) ToCadastrarPrestadorInput() (*input.CadastrarPrestadorInput, error) {
	voCPF, err := domain.NovoCPF(r.Cpf)
	if err != nil {
		return nil, err
	}

	// Gerar hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &input.CadastrarPrestadorInput{
		ID:          xid.New().String(),
		Nome:        r.Nome,
		CPF:         voCPF.String(),
		Email:       r.Email,
		Telefone:    r.Telefone,
		Senha:       string(hash), // Passa o hash ao invés da senha em texto
		ImagemUrl:   r.ImagemUrl,
		CatalogoIDs: r.CatalogoIDs,
	}, nil
}

func ValidaCPF(cpf string) (domain.CPF, error) {
	return domain.NovoCPF(cpf)
}
