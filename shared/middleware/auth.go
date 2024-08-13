package middleware

import (
	 // Importa tu paquete de configuración
	"log"

	"net/http"
	"strings"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, claims, err := ParseToken(tokenString)
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		role, exists := c.Get("role")
        log.Println("Role exists:", exists) // Verifica si el rol existe
        log.Println("Role value:", role)    // Imprime el valor del rol

        // Asegúrate de que el rol exista en el contexto
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found"})
            return
        }

        var roleInt64 int64
        switch v := role.(type) {
        case int64:
            roleInt64 = v
        case float64:
            roleInt64 = int64(v)
        case string:
            // Si el rol es un string, intenta convertirlo a int64 si es necesario
            // Para este ejemplo, se ignora el caso de string
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            return
        default:
            log.Println("Unsupported role type:", v)
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid role type"})
            return
        }

        if roleInt64 == 1 {
            log.Println("Admin access granted") // Informa si se concede acceso
            c.Next() // Si el rol es 1, permite el acceso
            return
        }

        log.Println("Access denied for role:", roleInt64) // Informa si se niega el acceso
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
    }
}

type Claims struct {
    UserID int64  `json:"id"`
	Username string `json:"username"`
    Role    int64 `json:"role"`
	Exp   int64  `json:"exp"`
		
    jwt.RegisteredClaims
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Utiliza la clave secreta desde las variables de entorno
        jwtSecret := []byte(os.Getenv("JWT_SECRET"))
        return jwtSecret, nil
    })
    return token, claims, err
}