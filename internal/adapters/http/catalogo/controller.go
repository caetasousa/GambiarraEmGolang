package catalogo

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"

	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatalogoController struct {
	criarCatalogoService *service.CatalogoService
}

func NovoCatalogoController(criarCatalogoService *service.CatalogoService) *CatalogoController {
	return &CatalogoController{criarCatalogoService: criarCatalogoService}
}

// @Summary Cria um novo catálogo de serviços
// @Description Cadastra um serviço que pode ser oferecido por um prestador
// @Tags Catalogos
// @Accept json
// @Produce json
// @Param catalogo body request_catalogo.CatalogoRequest true "Dados do Catálogo"
// @Success 201 {object} response_catalogo.CatalogoResponse "Catálogo criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Catálogo já existente"
// @Failure 500 {object} domain.ErrorResponse "Erro interno"
// @Router /catalogos [post]
func (ctl *CatalogoController) PostCatalogo(c *gin.Context) {
	var req request_catalogo.CatalogoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := req.ToCatalogoInput()

	out, err := ctl.criarCatalogoService.Cadastra(cmd)
	if err != nil {
		// Erros de validação do domínio retornam 400
		switch err {
		case domain.ErrNomeInvalido,
			domain.ErrDuracaoInvalida,
			domain.ErrPrecoInvalido,
			domain.ErrCategoriaInvalida:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			// Erros de infraestrutura ou inesperados retornam 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao criar catálogo"})
			return
		}
	}

	resp := response_catalogo.FromCatalogoResponse(*out)
	c.JSON(http.StatusCreated, resp)
}

// GetCatalogoPorID é o handler para a rota GET /catalogos/:id
// @Summary Busca um catálogo pelo ID
// @Description Retorna os dados de um catálogo específico usando seu ID.
// @Tags Catalogos
// @Accept json
// @Produce json
// @Param id path string true "ID do Catálogo"
// @Success 200 {object} response_catalogo.CatalogoResponse "Catálogo encontrado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "ID inválido fornecido"
// @Failure 404 {object} domain.ErrorResponse "Catálogo não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor ou falha de infraestrutura"
// @Router /catalogos/{id} [get]
func (ctl *CatalogoController) GetCatalogoPorID(c *gin.Context) {
	id := c.Param("id")

	catalogo, err := ctl.criarCatalogoService.BuscarPorId(id)
	if err != nil {
		switch err {
		case service.ErrCatalogoNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar catálogo"})
			return
		}
	}

	resp := response_catalogo.FromCatalogoResponse(*catalogo)
	c.JSON(http.StatusOK, resp)
}

// GetCatalogos godoc
// @Summary Lista todos os catálogos com paginação
// @Description Retorna uma lista de catálogos, com page e limit para paginação
// @Tags Catalogos
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Quantidade de itens por página" default(10)
// @Success 200 {object} response_catalogo.CatalogoListResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /catalogos [get]
func (ctl *CatalogoController) GetCatalogos(c *gin.Context) {
	var req request_catalogo.CatalogoListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	in := req.ToInputCatalogo()

	out, total, err := ctl.criarCatalogoService.Listar(in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao listar catálogos"})
		return
	}

	resp := response_catalogo.CatalogoListResponse{
		Data:  out,
		Page:  in.Page,
		Limit: in.Limit,
		Total: total,
	}

	c.JSON(http.StatusOK, resp)
}

// AtualizarCatalogo godoc
// @Summary Atualiza um catálogo existente
// @Description Atualiza os dados de um catálogo pelo ID
// @Tags Catalogos
// @Accept json
// @Produce json
// @Param id path string true "ID do Catálogo"
// @Param catalogo body request_catalogo.CatalogoUpdateRequest true "Dados atualizados do Catálogo"
// @Success 204 "Catálogo atualizado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 404 {object} domain.ErrorResponse "Catálogo não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /catalogos/{id} [put]
func (ctl *CatalogoController) Atualizar(c *gin.Context) {
	id := c.Param("id")

	var req request_catalogo.CatalogoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := req.ToCatalogoUpdateInput()
	input.ID = id

	err := ctl.criarCatalogoService.Atualizar(input)
	if err != nil {
		switch err {
		case service.ErrCatalogoNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case domain.ErrDuracaoInvalida,
			domain.ErrPrecoInvalido,
			domain.ErrNomeInvalido,
			domain.ErrCategoriaInvalida:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao atualizar catálogo"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// DeletarCatalogo godoc
// @Summary Deleta um catálogo existente
// @Description Remove um catálogo pelo ID
// @Tags Catalogos
// @Accept json
// @Produce json
// @Param id path string true "ID do Catálogo"
// @Success 204 "Catálogo deletado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "ID inválido"
// @Failure 404 {object} domain.ErrorResponse "Catálogo não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /catalogos/{id} [delete]
func (ctl *CatalogoController) Deletar(c *gin.Context) {
	id := c.Param("id")

	err := ctl.criarCatalogoService.Deletar(id)
	if err != nil {
		switch err {
		case service.ErrCatalogoNaoEncontrado:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case service.ErrFalhaInfraestrutura:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao deletar catálogo"})
			return
		}
	}

	c.Status(http.StatusNoContent)
}
