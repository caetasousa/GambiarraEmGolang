package teste

import (
	"encoding/json"
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria Silva",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img-original.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err := json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// Atualizar prestador
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "Maria Santos Atualizada",
		Email:       "maria.santos@email.com",
		Telefone:    "62988888888",
		ImagemUrl:   "https://exemplo.com/img-atualizada.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, createResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// Buscar e validar
	rrGet := SetupGetPrestadorRequest(router, createResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestadorAtualizado response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)
	require.NoError(t, err)

	assert.Equal(t, createResp.ID, prestadorAtualizado.ID)
	assert.Equal(t, "04423258196", prestadorAtualizado.Cpf)
	assert.Equal(t, "Maria Santos Atualizada", prestadorAtualizado.Nome)
	assert.Equal(t, "maria.santos@email.com", prestadorAtualizado.Email)
	assert.Equal(t, "62988888888", prestadorAtualizado.Telefone)
}

func TestUpdatePrestador_PrestadorNaoEncontrado(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	rr := SetupPutPrestadorRequest(router, "id-inexistente", updateInput)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "prestador não encontrado")
}

func TestUpdatePrestador_CatalogoInexistente(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")

	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva Atualizado",
		Email:       "joao.novo@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img2.jpg",
		CatalogoIDs: []string{"catalogo-inexistente"},
	}

	rr := SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "catálogo")
	require.Contains(t, rr.Body.String(), "não existe")
}

func TestUpdatePrestador_AtualizarCatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()

	// Criar dois catálogos
	catalogo1Input := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrCatalogo1 := SetupPostCatalogoRequest(router, catalogo1Input)
	require.Equal(t, http.StatusCreated, rrCatalogo1.Code)

	var catalogo1Resp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo1.Body.Bytes(), &catalogo1Resp)
	require.NoError(t, err)

	catalogo2Input := request_catalogo.CatalogoRequest{
		Nome:          "Barba",
		DuracaoPadrao: 20,
		Preco:         2000.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img2.jpg",
	}
	rrCatalogo2 := SetupPostCatalogoRequest(router, catalogo2Input)
	require.Equal(t, http.StatusCreated, rrCatalogo2.Code)

	var catalogo2Resp response_catalogo.CatalogoResponse
	err = json.Unmarshal(rrCatalogo2.Body.Bytes(), &catalogo2Resp)
	require.NoError(t, err)

	// Criar prestador com primeiro catálogo
	prestadorResp := CriarPrestadorValido(t, router, catalogo1Resp.ID, "04423258196")

	// Atualizar com ambos os catálogos
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogo1Resp.ID, catalogo2Resp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// Validar
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestadorAtualizado response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)
	require.NoError(t, err)

	assert.Len(t, prestadorAtualizado.Catalogo, 2)
}

func TestUpdatePrestador_DadosInvalidos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")

	testCases := []struct {
		name        string
		input       request_prestador.PrestadorUpdateRequest
		expectedMsg string
	}{
		{
			name: "Nome muito curto",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "Jo",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Nome",
		},
		{
			name: "Email inválido",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "email-invalido",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Email",
		},
		{
			name: "Telefone muito curto",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "joao@email.com",
				Telefone:    "1234567",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Telefone",
		},
		{
			name: "URL inválida",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "url-invalida",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "ImagemUrl",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := SetupPutPrestadorRequest(router, prestadorResp.ID, tc.input)
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Contains(t, rr.Body.String(), tc.expectedMsg)
		})
	}
}

func TestUpdatePrestador_SemCatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")

	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{},
	}

	rr := SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rr.Code)
}
func TestUpdatePrestador_CPFNaoMuda(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")

	cpfOriginal := prestadorResp.Cpf

	// Tentar atualizar (CPF não está no request de update, mas vamos validar)
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "Nome Atualizado",
		Email:       "novo@email.com",
		Telefone:    "62988888888",
		ImagemUrl:   "https://exemplo.com/img.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// Verificar que CPF permanece o mesmo
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	var prestadorAtualizado response_prestador.PrestadorResponse
	json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)

	assert.Equal(t, cpfOriginal, prestadorAtualizado.Cpf, "CPF não deve mudar")
}

// 5. Verificar que agenda é mantida após atualização
func TestUpdatePrestador_MantémAgenda(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	// Criar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	// Atualizar prestador
	catalogoResp := CriarCatalogoValido(t, router)
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "Nome Atualizado",
		Email:       "novo@email.com",
		Telefone:    "62988888888",
		ImagemUrl:   "https://exemplo.com/img.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)

	// Verificar que agenda ainda existe
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	var prestadorAtualizado response_prestador.PrestadorResponse
	json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)

	assert.NotEmpty(t, prestadorAtualizado.Agenda, "Agenda deve ser mantida")
	assert.Len(t, prestadorAtualizado.Agenda, 1)
	assert.Equal(t, "2030-01-03", prestadorAtualizado.Agenda[0].Data)
}
