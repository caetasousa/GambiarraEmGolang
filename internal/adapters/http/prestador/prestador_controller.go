package prestador

import (
	"meu-servico-agenda/internal/adapters/http/prestador/request"
	"meu-servico-agenda/internal/core/application/service"
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
// PostPrestador é o handler para a rota POST /prestador
// @Summary Cadastra um novo prestadores
// @Description Recebe os dados necessários para registrar um novo prestador.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param prestador body request.PrestadorRequest true "Dados do Prestador"
// @Success 201 {object} domain.Prestador "Prestador criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação do binding)"
// @Failure 409 {object} domain.ErrorResponse "Prestador já cadastrado ou conflito de dados"
// @Failure 500 {object} domain.ErrorResponse "Falha na persistência de dados ou erro interno"
// @Router /prestadores [post]
func (prc *PrestadorController) PostPrestador(c *gin.Context) {
	var input request.PrestadorRequest

	// 1️⃣ Validação de binding do JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// 2️⃣ Orquestração pelo service
	prestador, err := prc.prestadorService.Cadastra(&input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// 3️⃣ Retorno HTTP 201 com o prestador criado
	c.JSON(http.StatusCreated, prestador)
}

// @Summary Define a agenda diária de um prestador
// @Description Adiciona uma nova agenda diária à disponibilidade do prestador.
// @Tags Prestadores
// @Accept json
// @Produce json
// @Param id path string true "ID do prestador"
// @Param agenda body request.AgendaDiariaRequest true "Agenda diária"
// @Success 204 "Agenda adicionada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Conflito de agenda"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /prestadores/{id}/agenda [put]
func (prc *PrestadorController) PutAgenda(c *gin.Context) {
	prestadorID := c.Param("id")

	var input request.AgendaDiariaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := prc.prestadorService.AdicionarAgenda(prestadorID, &input); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
