package prestador

import (
	"errors"

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

// @Summary Define a agenda diária de um prestador
// @Description Adiciona uma nova agenda diária à disponibilidade do prestador.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param id path string true "ID do prestador"
// @Param agenda body request_prestador.AgendaDiariaRequest true "Agenda diária"
// @Success 204 "Agenda adicionada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Conflito de agenda"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id}/agenda [put]
func (prc *PrestadorController) PutAgenda(c *gin.Context) {
	prestadorID := c.Param("id")

	var input request_prestador.AgendaDiariaRequest

	// 1️⃣ Binding / validação de formato
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Adapter → Command (mantém time.Parse e validações)
	cmd, err := input.ToAdicionarAgendaInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3️⃣ Caso de uso
	err = prc.prestadorService.AdicionarAgenda(prestadorID, cmd)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPrestadorNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, domain.ErrAgendaDuplicada),
			errors.Is(err, service.ErrPrestadorInativo),
			errors.Is(err, domain.ErrPrestadorInativo):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case errors.Is(err, domain.ErrAgendaSemIntervalos),
			errors.Is(err, domain.ErrIntervaloHorarioInvalido),
			errors.Is(err, domain.ErrDataEstaNoPassado):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno"})
		}
		return
	}

	c.Status(http.StatusNoContent)
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
// @Summary Lista todos os prestadores
// @Description Retorna uma lista paginada de prestadores com seus catálogos e agendas
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)" default(1) minimum(1)
// @Param limit query int false "Itens por página (padrão: 10, máximo: 100)" default(10) minimum(1) maximum(100)
// @Success 200 {object} response_prestador.PrestadorListResponse "Lista de prestadores retornada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Parâmetros de paginação inválidos"
// @Failure 500 {object} domain.ErrorResponse "Erro interno ao listar prestadores"
// @Router /prestadores [get]
func (prc *PrestadorController) GetPreestadores(c *gin.Context) {
	var req request_prestador.PrestadorListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	in := req.ToInputPrestador()

	out, total, err := prc.prestadorService.Listar(in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao listar catálogos"})
		return
	}

	resp := response_prestador.PrestadorListResponse{
		Data:  out,
		Page:  in.Page,
		Limit: in.Limit,
		Total: total,
	}

	c.JSON(http.StatusOK, resp)
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