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
