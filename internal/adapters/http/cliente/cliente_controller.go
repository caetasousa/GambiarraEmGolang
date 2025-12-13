package cliente

import (
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/core/application/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ClienteController struct {
	novoCliente *service.ServiceCliente
}

func NovoClienteController(novoCliente *service.ServiceCliente) *ClienteController {
	return &ClienteController{novoCliente: novoCliente}
}

// PostCliente é o handler para a rota POST /cliente
// @Summary Cadastra um novo cliente
// @Description Recebe dados de nome, email e telefone para registrar um novo cliente.
// @Tags Clientes
// @Accept json
// @Produce json
// @Param cliente body request.ClienteRequest true "Dados do Cliente"
// @Success 201 {object} domain.Cliente "Cliente criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação do binding)"
// @Failure 409 {object} domain.ErrorResponse "Cliente já cadastrado (Email ou Telefone já existe)"
// @Failure 500 {object} domain.ErrorResponse "Falha na persistência de dados ou erro interno"
// @Router /clientes [post]
func (ctrl *ClienteController) PostCliente(c *gin.Context) {
	var input request.ClienteRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Cria domínio e valida regras de negócio
	clienteDomain, err := input.ToCliente()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Persiste usando service
	cliente, err := ctrl.novoCliente.Cadastra(clienteDomain)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}

// @Summary Busca um cliente pelo ID
// @Description Retorna os dados de um cliente específico usando seu ID.
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path string true "ID do Cliente"
// @Success 200 {object} domain.Cliente "Cliente encontrado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "ID inválido fornecido (ex: formato incorreto se houver validação de formato de ID)"
// @Failure 404 {object} domain.ErrorResponse "Cliente não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor ou falha de infraestrutura"
// @Router /clientes/{id} [get]
func (ctrl *ClienteController) GetCliente(c *gin.Context) {
	id := c.Param("id")

	cliente, err := ctrl.novoCliente.BuscarPorId(id)

	if err != nil {
		errorMessage := err.Error()

		// 1. TRATAMENTO DO 404: Se o erro for a mensagem específica de "não encontrado"
		if errorMessage == "cliente não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
			return
		}

		// 2. TRATAMENTO DO 500: Qualquer outro erro é tratado como falha de infraestrutura
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar cliente: " + errorMessage})
		return
	}

	// Se o serviço não retorna erro, mas retorna nil (caso o serviço seja simplificado)
	if cliente == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}
