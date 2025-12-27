package teste

import (
	"encoding/json"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"
	"net/http"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fluxo completo: criar -> adicionar agenda -> atualizar -> buscar
func TestFluxoCompleto_CriarAgendarAtualizar(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// 1. Criar prestador
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")

	// 2. Adicionar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// 3. Atualizar prestador
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "Nome Atualizado",
		Email:       "atualizado@email.com",
		Telefone:    "62999999999",
		ImagemUrl:   "https://exemplo.com/new.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, prestadorResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// 4. Buscar e validar TUDO
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	var final response_prestador.PrestadorResponse
	json.Unmarshal(rrGet.Body.Bytes(), &final)

	assert.Equal(t, "Nome Atualizado", final.Nome)
	assert.Equal(t, "04423258196", final.Cpf) // CPF não mudou
	assert.Len(t, final.Agenda, 1)            // Agenda mantida
	assert.Len(t, final.Catalogo, 1)          // Catálogo mantido
}

// Desativar prestador deve impedir nova agenda
func TestPrestadorInativo_NaoPermiteNovaAgenda(t *testing.T) {
	router, prestadorResp, repo := CriarPrestadorValidoParaTeste(t)

	// Desativar prestador
	prestadorResp.Ativo = false
	repo.Salvar(&prestadorResp)

	// Tentar adicionar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusConflict, rr.Code)
}

// 4. Concorrência - dois requests simultâneos
func TestConcorrencia_CriacaoSimultanea(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	// Criar 2 requests simultâneos com mesmo CPF
	var wg sync.WaitGroup
	results := make([]int, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			rr := SetupPostPrestadorRequest(router, prestadorInput)
			results[index] = rr.Code
		}(i)
	}

	wg.Wait()

	// Um deve ser 201 e outro 409 (CPF duplicado)
	codes := []int{results[0], results[1]}
	sort.Ints(codes)
	assert.Equal(t, []int{201, 409}, codes)
}
