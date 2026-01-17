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

// PostCliente é o handler para a rota POST /clientes
// @Summary Cadastra um novo cliente
// @Description Recebe dados de nome, email, telefone e senha para registrar um novo cliente.
// @Tags Clientes
// @Accept json
// @Produce json
// @Param cliente body request.ClienteRequest true "Dados do Cliente"
// @Success 201 {object} output.BuscarClienteOutput "Cliente criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação do binding)"
// @Failure 409 {object} domain.ErrorResponse "Cliente já cadastrado (Email já existe)"
// @Failure 500 {object} domain.ErrorResponse "Falha na persistência de dados ou erro interno"
// @Router /clientes [post]
func (ctrl *ClienteController) PostCliente(c *gin.Context) {
	var req request.ClienteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Converter para input
	inputData, err := req.ToCadastrarClienteInput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}

	// Cadastrar usando service
	cliente, err := ctrl.novoCliente.Cadastra(inputData)
	if err != nil {
		// Tratamento específico de erros
		switch err.Error() {
		case "email já cadastrado":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "email já está em uso por um prestador":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "senha deve ter no mínimo 8 caracteres":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar cliente"})
		}
		return
	}

	c.JSON(http.StatusCreated, cliente)
}

// @Summary Busca um cliente pelo ID
// @Description Retorna os dados de um cliente específico usando seu ID. Requer autenticação JWT.
// @Tags Clientes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do Cliente"
// @Success 200 {object} output.BuscarClienteOutput "Cliente encontrado com sucesso"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - você só pode acessar seus próprios dados"
// @Failure 404 {object} domain.ErrorResponse "Cliente não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor ou falha de infraestrutura"
// @Router /clientes/{id} [get]
func (ctrl *ClienteController) GetCliente(c *gin.Context) {
	id := c.Param("id")

	cliente, err := ctrl.novoCliente.BuscarPorId(id)

	if err != nil {
		errorMessage := err.Error()

		// TRATAMENTO DO 404
		if errorMessage == "cliente não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cliente não encontrado"})
			return
		}

		// TRATAMENTO DO 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao buscar cliente"})
		return
	}

	c.JSON(http.StatusOK, cliente)
}
