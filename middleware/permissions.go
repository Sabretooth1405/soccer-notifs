package middleware

import (
	"fmt"
	"net/http"
	"os"
	"soccer-notifs/initializers"
	"soccer-notifs/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *gin.Context) {
	tokenString := c.GetHeader("Token")
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SIGNINING_KEY")), nil
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("b")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userId := claims["sub"]
		var user models.User
		result := initializers.DB.First(&user, userId)
		if result.Error != nil {
			fmt.Println("c")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user",user)
	} else {
		fmt.Println("d")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
