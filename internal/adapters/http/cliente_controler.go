package Http

import (
	"meu-servico-agenda/internal/core/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClienteController struct {
	
	cadastrador *services.CadastroDeCliente
}

func NovoClienteController(cadastrador *services.CadastroDeCliente) *ClienteController {
	return &ClienteController{cadastrador: cadastrador}
}

func (ctrl *ClienteController) PostCliente(c *gin.Context) {
	var input services.CadastrarClienteInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	cliente, err := ctrl.cadastrador.Executar(input)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cliente)
}
