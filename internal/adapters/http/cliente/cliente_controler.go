package Http

import (
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/core/application/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ClienteController struct {
	novoCliente *services.CadastroDeCliente
}

func NovoClienteController(novoCliente *services.CadastroDeCliente) *ClienteController {
	return &ClienteController{novoCliente: novoCliente}
}

// PostCliente é o handler para a rota POST /clientes
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

	cliente, err := ctrl.novoCliente.Executar(input)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}
