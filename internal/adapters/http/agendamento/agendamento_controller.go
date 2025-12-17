package agendamento

import (
	requesta_gendamento "meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/core/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgendamentoController struct {
	agendamentoService *service.AgendamentoService
}

func NovoAgendamentoController(ag *service.AgendamentoService) *AgendamentoController {
	return &AgendamentoController{
		agendamentoService: ag,
	}
}

func (ag *AgendamentoController) PostAgendamento(c *gin.Context) {
	// 1️⃣ Recebe os dados da requisição no formato JSON
	var input requesta_gendamento.AgendamentoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		// Se houver erro ao fazer o bind, retorna erro 400 (Bad Request)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// 2️⃣ Chama o serviço para cadastrar o agendamento
	// Agora passamos o 'input' diretamente para o serviço
	agendamento, err := ag.agendamentoService.CadastraAgendamento(input)

	if err != nil {
		// Se ocorrer algum erro ao cadastrar o agendamento, retornamos uma resposta de erro
		// Exemplo de retorno de erro: caso o horário esteja indisponível
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// 3️⃣ Se tudo correr bem, retornamos o agendamento criado
	c.JSON(http.StatusCreated, agendamento)
}
