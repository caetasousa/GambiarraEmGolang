package prestador

import (
	"errors"
	"fmt"
	"time"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"

	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrestadorController struct {
	prestadorService *service.PrestadorService
}

func NovoPrestadorController(ps *service.PrestadorService) *PrestadorController {
	return &PrestadorController{
		prestadorService: ps,
	}
}

// @Summary Cadastra um novo prestador
// @Description Recebe os dados necessários para registrar um novo prestador.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param prestador body request_prestador.PrestadorRequest true "Dados do Prestador"
// @Success 201 {object} response_prestador.PrestadorPostResponse "Prestador criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação do binding)"
// @Failure 409 {object} domain.ErrorResponse "Prestador já cadastrado ou conflito de dados"
// @Failure 500 {object} domain.ErrorResponse "Falha na persistência de dados ou erro interno"
// @Router /prestadores [post]
func (prc *PrestadorController) PostPrestador(c *gin.Context) {
	var input request_prestador.PrestadorRequest

	// 1️⃣ Validação de binding / formato (continua igual)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// 2️⃣ Adapter → Command (mantém validações existentes)
	cmd, err := input.ToCadastrarPrestadorInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3️⃣ Chamada do caso de uso
	prestador, err := prc.prestadorService.Cadastra(cmd)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCPFJaCadastrado):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case errors.Is(err, service.ErrCatalogoNaoExiste):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 4️⃣ Response
	resp := response_prestador.FromCriarPrestadorOutput(*prestador)
	c.JSON(http.StatusCreated, resp)
}

// @Summary Consulta prestador pelo ID
// @Description Retorna informações do prestador, incluindo catálogo de serviços e agenda diária.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param id path string true "ID do prestador"
// @Success 200 {object} response_prestador.PrestadorResponse "Prestador encontrado com sucesso"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id} [get]
func (prc *PrestadorController) GetPrestador(c *gin.Context) {
	id := c.Param("id")

	out, err := prc.prestadorService.BuscarPorId(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPrestadorNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "erro interno",
			})
			return
		}
	}

	resp := response_prestador.FromPrestadorOutput(*out)
	c.JSON(http.StatusOK, resp)
}

// UpdatePrestador godoc
// @Summary Atualiza um prestador existente
// @Description Atualiza os dados cadastrais de um prestador, incluindo nome, email, telefone, imagem e catálogos de serviços associados. O CPF não pode ser alterado.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param id path string true "ID do prestador"
// @Param prestador body request_prestador.PrestadorUpdateRequest true "Dados atualizados do prestador"
// @Success 204 "Prestador atualizado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou catálogo não existe"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /prestadores/{id} [put]
func (prc *PrestadorController) UpdatePrestador(c *gin.Context) {
	id := c.Param("id")

	var req request_prestador.PrestadorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// ✅ Já captura erros de validação (nome, email, telefone, url)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := req.ToAlterarPrestadorInput()
	input.Id = id

	err := prc.prestadorService.Atualizar(input)
	if err != nil {
		switch err {
		case service.ErrPrestadorNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrCatalogoNaoExiste:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case domain.ErrPrestadorDeveTerCatalogo: // ✅ NOVO
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao atualizar prestador"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// GetPrestadores godoc
// @Summary Lista prestadores filtrados por status
// @Description Retorna lista paginada de prestadores ativos ou inativos
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param limit query int false "Itens por página (padrão: 10, máximo: 100)"
// @Param ativo query boolean true "Status do prestador (obrigatório)"
// @Success 200 {object} response_prestador.PrestadorListResponse
// @Failure 400 {object} domain.ErrorResponse "Parâmetro 'ativo' é obrigatório"
// @Router /prestadores [get]
func (prc *PrestadorController) GetPrestadores(c *gin.Context) {
	var req request_prestador.PrestadorListRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Parâmetro 'ativo' é obrigatório: %s", err.Error()),
		})
		return
	}

	input := req.ToInputPrestador()

	prestadores, total, err := prc.prestadorService.ListarPrestadores(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao listar prestadores",
		})
		return
	}

	response := response_prestador.PrestadorListResponse{
		Data:  prestadores,
		Page:  input.Page,
		Limit: input.Limit,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Inativa um prestador
// @Description Inativa um prestador, impedindo que ele receba novos agendamentos
// @Tags Prestadores
// @Produce json
// @Param id path string true "ID do prestador"
// @Success 204 "Prestador inativado com sucesso"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id}/inativar [put]
func (prc *PrestadorController) InativarPrestador(c *gin.Context) {
	id := c.Param("id")

	err := prc.prestadorService.Inativar(id)
	if err != nil {
		switch err {
		case service.ErrPrestadorNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @Summary Ativa um prestador
// @Description Ativa um prestador, permitindo que ele receba novos agendamentos
// @Tags Prestadores
// @Produce json
// @Param id path string true "ID do prestador"
// @Success 204 "Prestador ativado com sucesso"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id}/ativar [put]
func (prc *PrestadorController) AtivarPrestador(c *gin.Context) {
	id := c.Param("id")

	err := prc.prestadorService.Ativar(id)
	if err != nil {
		switch err {
		case service.ErrPrestadorNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// @Summary Cria ou atualiza uma agenda
// @Description Cria uma nova agenda ou atualiza uma existente para a data especificada
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param id path string true "ID do prestador"
// @Param agenda body request_prestador.AgendaDiariaRequest true "Dados da agenda"
// @Success 200 "Agenda criada ou atualizada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Prestador inativo"
// @Router /prestadores/{id}/agenda [put]
func (prc *PrestadorController) PutAgenda(c *gin.Context) {
	prestadorID := c.Param("id")

	var input request_prestador.AgendaDiariaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Dados inválidos: %s", err.Error()),
		})
		return
	}

	cmd, err := input.ToAdicionarAgendaInput() // Um método só
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.PrestadorID = prestadorID

	err = prc.prestadorService.SalvarAgenda(cmd)
	if err != nil {
		switch err {
		case service.ErrPrestadorNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrPrestadorInativo:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case domain.ErrIntervaloHorarioInvalido,
			domain.ErrAgendaSemIntervalos,
			domain.ErrDataEstaNoPassado,
			domain.ErrIntervalosSesobrepoe: // ✅ NOVO
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao salvar agenda",
			})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// DeleteAgenda godoc
// @Summary Deleta uma agenda
// @Description Remove uma agenda de um prestador em uma data específica
// @Tags Prestadores
// @Produce json
// @Param id path string true "ID do prestador"
// @Param data query string true "Data da agenda (formato: 2006-01-02)"
// @Success 204 "Agenda deletada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Data inválida"
// @Failure 404 {object} domain.ErrorResponse "Prestador ou agenda não encontrada"
// @Failure 409 {object} domain.ErrorResponse "Prestador inativo"
// @Router /prestadores/{id}/agenda [delete]
func (prc *PrestadorController) DeleteAgenda(c *gin.Context) {
	prestadorID := c.Param("id")
	data := c.Query("data") // Ex: ?data=2030-01-03

	// Validar formato da data
	_, err := time.Parse("2006-01-02", data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data inválida. Use o formato YYYY-MM-DD",
		})
		return
	}

	err = prc.prestadorService.DeletarAgenda(prestadorID, data)
	if err != nil {
		switch err {
		case service.ErrPrestadorNaoEncontrado,
			service.ErrAgendaNaoEncontrada:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrPrestadorInativo:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao deletar agenda",
			})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// GetPrestadoresPorData godoc
// @Summary Lista prestadores disponíveis em uma data específica
// @Description Retorna lista paginada de prestadores ativos que possuem agenda configurada para a data informada
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param data query string true "Data da disponibilidade (formato: YYYY-MM-DD)"
// @Param page query int false "Número da página (padrão: 1)"
// @Param limit query int false "Itens por página (padrão: 10, máximo: 100)"
// @Success 200 {object} response_prestador.PrestadorListResponse "Lista de prestadores disponíveis"
// @Failure 400 {object} domain.ErrorResponse "Parâmetro 'data' é obrigatório ou formato inválido"
// @Failure 500 {object} domain.ErrorResponse "Erro interno ao buscar prestadores"
// @Router /prestadores/disponiveis [get]
func (prc *PrestadorController) GetPrestadoresPorData(c *gin.Context) {
	var req request_prestador.BuscarPrestadoresDataRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Parâmetro 'data' é obrigatório: %s", err.Error()),
		})
		return
	}

	input := req.ToInputPrestador()

	prestadores, total, err := prc.prestadorService.BuscarPrestadoresDisponiveisPorData(input)
	if err != nil {
		switch err {
		case service.ErrAoBuscarPrestadoresDisponiveis,
			service.ErrAoContarPrestadoresDisponiveis:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrFormatoDataInvalido, domain.ErrDataEstaNoPassado:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao listar prestadores",
			})
			return
		}
	}

	response := response_prestador.PrestadorListResponse{
		Data:  prestadores,
		Page:  input.Page,
		Limit: input.Limit,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}
