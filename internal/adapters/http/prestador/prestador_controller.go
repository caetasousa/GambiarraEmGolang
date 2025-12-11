package prestador

import (
	"meu-servico-agenda/internal/adapters/http/prestador/request"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrestadorController struct {
	criarPrestadorService *service.PrestadorService
	catalogoRepo          port.CatalogoRepositorio
}

func NovoPrestadorController(criarPrestadorService *service.PrestadorService, catalogoRepo port.CatalogoRepositorio) *PrestadorController {
	return &PrestadorController{criarPrestadorService: criarPrestadorService, catalogoRepo: catalogoRepo}
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

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	prestador, err := input.ToPrestador(prc.catalogoRepo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prestador, err = prc.criarPrestadorService.Cadastra(prestador)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prestador)
}
