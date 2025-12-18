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
// @Success 201 {object} response.PrestadorPostResponse "Prestador criado com sucesso"
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
	cmd, err := input.ToCommand()
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
	resp := response_prestador.FromPostPrestador(prestador)
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
	cmd, err := input.ToCommand()
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
			errors.Is(err, domain.ErrPrestadorInativo):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case errors.Is(err, domain.ErrAgendaSemIntervalos),
			errors.Is(err, domain.ErrIntervaloHorarioInvalido):
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
// @Success 200 {object} response.PrestadorResponse "Prestador encontrado com sucesso"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id} [get]
func (prc *PrestadorController) GetPrestador(c *gin.Context) {
	id := c.Param("id")

	prestador, err := prc.prestadorService.BuscarPorId(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	resp := response_prestador.FromPrestador(prestador)

	c.JSON(http.StatusOK, resp)
}
