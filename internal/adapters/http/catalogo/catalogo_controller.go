package catalogo

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"

	"meu-servico-agenda/internal/core/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatalogoController struct {
	criarCatalogoService *service.CatalogoService
}

func NovoCatalogoController(criarCatalogoService *service.CatalogoService) *CatalogoController {
	return &CatalogoController{criarCatalogoService: criarCatalogoService}
}

// PostPrestador é o handler para a rota POST /catalogos
// @Summary Cadastra um novo catálogo
// @Description Recebe os dados necessários para registrar um novo serviço no catálogo.
// @Tags Catálogos
// @Accept json
// @Produce json
// @Param catalogo body request.CatalogoRequest true "Dados do Catálogo"
// @Success 201 {object} response.CatalogoResponse "Catálogo criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação do binding)"
// @Failure 409 {object} domain.ErrorResponse "Catálogo já cadastrado ou conflito de dados"
// @Failure 500 {object} domain.ErrorResponse "Falha na persistência de dados ou erro interno"
// @Router /catalogos [post]
func (ctl *CatalogoController) PostCatalogo(c *gin.Context) {
	var input request_catalogo.CatalogoRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := input.ToCommand()

	catalogoSalvo, err := ctl.criarCatalogoService.Cadastra(cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := response_catalogo.FromCatalogo(catalogoSalvo)
	c.JSON(http.StatusCreated, resp)
}

// GetCatalogoPorID é o handler para a rota GET /catalogos/:id
// @Summary Busca um catálogo pelo ID
// @Description Retorna os dados de um catálogo específico usando seu ID.
// @Tags Catálogos
// @Accept json
// @Produce json
// @Param id path string true "ID do Catálogo"
// @Success 200 {object} response.CatalogoResponse "Catálogo encontrado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "ID inválido fornecido"
// @Failure 404 {object} domain.ErrorResponse "Catálogo não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor ou falha de infraestrutura"
// @Router /catalogos/{id} [get]
func (ctl *CatalogoController) GetCatalogoPorID(c *gin.Context) {
	id := c.Param("id")

	catalogo, err := ctl.criarCatalogoService.BuscarPorId(id)

	if err != nil {
		if err.Error() == "catalogo nao encontrado" { // compara string exata
			c.JSON(http.StatusNotFound, gin.H{"error": "Catalogo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := response_catalogo.FromCatalogo(catalogo)

	c.JSON(http.StatusOK, resp)
}
