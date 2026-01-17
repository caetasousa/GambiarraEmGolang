package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/klassmann/cpfcnpj"
)

var ErrCPFInvalido = errors.New("cpf inválido")

type CPF struct {
	valor string
}

func NovoCPF(v string) (CPF, error) {
	limpo := cpfcnpj.Clean(v)
	if !cpfcnpj.ValidateCPF(limpo) || todosDigitosIguais(limpo) {
		return CPF{}, ErrCPFInvalido
	}
	return CPF{valor: limpo}, nil
}

func todosDigitosIguais(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func (c CPF) String() string {
	return c.valor
}

func (c CPF) Formatar() string {
	if len(c.valor) != 11 {
		return c.valor
	}
	return strings.Join([]string{
		c.valor[:3], ".",
		c.valor[3:6], ".",
		c.valor[6:9], "-",
		c.valor[9:],
	}, "")
}

// Scan implementa a interface sql.Scanner
func (c *CPF) Scan(value interface{}) error {
	if value == nil {
		c.valor = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("tipo inválido para CPF: %T", value)
	}
	c.valor = s
	return nil
}

// Value implementa a interface driver.Valuer
func (c CPF) Value() (driver.Value, error) {
	return c.valor, nil
}
