package request_prestador

import (
	"errors"
	"meu-servico-agenda/internal/core/domain"

	"github.com/klassmann/cpfcnpj"
)

type PrestadorRequest struct {
	Nome        string   `json:"nome" binding:"required,min=3,max=100" example:"joao" swagger:"desc('Nome do prestador')"`
	Cpf         string   `json:"cpf" binding:"required,len=11,numeric" example:"04423258196"`
	Email       string   `json:"email" binding:"omitempty,email" example:"joao@email.com" swagger:"desc('Email do prestador')"`
	Telefone    string   `json:"telefone" binding:"required,min=8,max=15" example:"62999677481" swagger:"desc('Telefone do prestador')"`
	CatalogoIDs []string `json:"catalogo_ids" binding:"omitempty,dive,required" swagger:"desc('IDs dos serviços no catálogo oferecidos pelo prestador')"`
}

func (r *PrestadorRequest) ToPrestador(catalogos []domain.Catalogo) (*domain.Prestador, error) {
	cpf, err := ValidaCPF(r.Cpf)
	if err != nil {
		return nil, err
	}
	return domain.NovoPrestador(r.Nome, cpf, r.Email, r.Telefone, catalogos)
}

func ValidaCPF(cpf string) (string, error) {
	if !cpfcnpj.ValidateCPF(cpf) {
		return "", errors.New("cpf inválido")
	}
	return cpf, nil
}
