package teste

import (
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"

	"github.com/stretchr/testify/require"
)

func TestPostPrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Email:       "joao@email.com",
		Cpf:         "04423258196",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code)
}

func TestPostPrestador_FalhaCatalogoInexistente(t *testing.T) {
	router, _ := SetupPostPrestador()

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{"catalogo-inexistente"},
	}

	rr := SetupPostPrestadorRequest(router, prestadorInput)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	require.Contains(t, rr.Body.String(), "não existe")
}

func TestPostPrestador_CPFDuplicado(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	ccCreate := SetupPostPrestadorRequest(router, prestadorInput)
	
	require.Equal(t, http.StatusCreated, rrCreate.Code)
	require.Equal(t, http.StatusConflict, ccCreate.Code)
}
func TestPostPrestador_CamposObrigatorios(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	testCases := []struct {
		name        string
		input       request_prestador.PrestadorRequest
		expectedMsg string
	}{
		{
			name: "Nome vazio",
			input: request_prestador.PrestadorRequest{
				Nome:        "",
				Cpf:         "04423258196",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Nome",
		},
		{
			name: "CPF inválido",
			input: request_prestador.PrestadorRequest{
				Nome:        "João Silva",
				Cpf:         "12345678900",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "cpf",
		},
		{
			name: "Email inválido",
			input: request_prestador.PrestadorRequest{
				Nome:        "João Silva",
				Cpf:         "04423258196",
				Email:       "email-invalido",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Email",
		},
		{
			name: "Telefone muito curto",
			input: request_prestador.PrestadorRequest{
				Nome:        "João Silva",
				Cpf:         "04423258196",
				Email:       "joao@email.com",
				Telefone:    "123",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Telefone",
		},
		{
			name: "Sem catálogos",
			input: request_prestador.PrestadorRequest{
				Nome:        "João Silva",
				Cpf:         "04423258196",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{},
			},
			expectedMsg: "CatalogoIDs", // ✅ Mensagem real do binding
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := SetupPostPrestadorRequest(router, tc.input)
			require.Equal(t, http.StatusBadRequest, rr.Code, "Response: %s", rr.Body.String())
			require.Contains(t, rr.Body.String(), tc.expectedMsg)
		})
	}
}