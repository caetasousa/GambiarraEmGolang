package middleware

import (
	"net/http"
	"strings"

	"meu-servico-agenda/internal/core/domain"
	"meu-servico-agenda/internal/infra/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida o token JWT e injeta os claims no contexto
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Pegar token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autenticação não fornecido"})
			c.Abort()
			return
		}

		// 2. Validar formato: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido. Use: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Validar e decodificar JWT
		claims, err := auth.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 4. Injetar claims no contexto Gin para próximos handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", string(claims.Role))

		c.Next()
	}
}

// RequireRole valida se o usuário tem uma das roles permitidas
func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			c.Abort()
			return
		}

		userRoleStr := userRole.(string)

		// Verificar se role está permitida
		for _, role := range allowedRoles {
			if userRoleStr == string(role) {
				c.Next()
				return
			}
		}

		// Se não tem permissão
		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado: permissão insuficiente"})
		c.Abort()
	}
}

// RequireOwnerOrAdmin valida se é o próprio usuário ou admin
// Usa o parâmetro "id" da rota para comparar com o user_id do token
func RequireOwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			c.Abort()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			c.Abort()
			return
		}

		resourceID := c.Param("id")
		userIDStr := userID.(string)
		userRoleStr := userRole.(string)

		// Admin pode tudo
		if userRoleStr == string(domain.RoleAdmin) {
			c.Next()
			return
		}

		// Usuário pode acessar apenas seus próprios dados
		if userIDStr == resourceID {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado: você só pode acessar seus próprios dados"})
		c.Abort()
	}
}

// RequireOwnerOrRole valida se é o próprio usuário ou tem uma das roles especificadas
// Útil para casos como: prestador pode editar própria agenda, admin também pode
func RequireOwnerOrRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			c.Abort()
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			c.Abort()
			return
		}

		resourceID := c.Param("id")
		userIDStr := userID.(string)
		userRoleStr := userRole.(string)

		// Verificar se tem uma das roles permitidas
		for _, role := range allowedRoles {
			if userRoleStr == string(role) {
				c.Next()
				return
			}
		}

		// Se não tem role permitida, verificar se é o próprio usuário
		if userIDStr == resourceID {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "acesso negado: permissão insuficiente"})
		c.Abort()
	}
}
