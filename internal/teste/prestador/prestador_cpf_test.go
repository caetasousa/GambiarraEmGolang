package teste

import (
	"testing"

	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

func TestNovoCPF(t *testing.T) {
	tests := []struct {
		name    string
		valor   string
		wantErr bool
	}{
		{"CPF Válido com pontos", "044.232.581-96", false},
		{"CPF Válido limpo", "04423258196", false},
		{"CPF Inválido curto", "123", true},
		{"CPF Inválido longo", "123456789012", true},
		{"CPF Inválido dígitos", "11111111111", true},
		{"Vazio", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NovoCPF(tt.valor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, domain.ErrCPFInvalido, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCPF_String(t *testing.T) {
	cpf, _ := domain.NovoCPF("044.232.581-96")
	assert.Equal(t, "04423258196", cpf.String())
}

func TestCPF_Formatar(t *testing.T) {
	cpf, _ := domain.NovoCPF("04423258196")
	assert.Equal(t, "044.232.581-96", cpf.Formatar())
}

func TestCPF_ScannerValuer(t *testing.T) {
	var c domain.CPF
	err := c.Scan("04423258196")
	assert.NoError(t, err)
	assert.Equal(t, "04423258196", c.String())

	val, err := c.Value()
	assert.NoError(t, err)
	assert.Equal(t, "04423258196", val)

	err = c.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, "", c.String())
}
