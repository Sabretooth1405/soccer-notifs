package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"soccer-notifs/initializers"
	"soccer-notifs/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	Timezone          string `json:"timezone,omitempty"`
	Email_is_verified bool   `json:"email_is_verified,omitempty"`
}

func Register(c *gin.Context) {
	var user User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "Failed to read body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password")
	}

	var result *gorm.DB
	var createUser models.User
	fmt.Println("Timezone:",user.Timezone != "" )
	if user.Timezone == "" {
		createUser = models.User{Email: user.Email, Password: string(hash)}
		result = initializers.DB.Create(&createUser)
	} else {
		createUser = models.User{Email: user.Email, Password: string(hash), Timezone: user.Timezone}
		result = initializers.DB.Create(&createUser)
	}
	if result.Error != nil {
		fmt.Println(result.Error)
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{
				"respone": "Email Already Registered",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"respone": "Email Already Registered",
			})
			return
		}
	}
	createUser.Password = "********"
	c.JSON(200, gin.H{
		"response": &createUser,
	})
}
func Login(c *gin.Context){
	var user User

	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "Failed to read body",
		})
		return
	}
	var findUser models.User
	result:=initializers.DB.Where("email = ? ",user.Email).First(&findUser)
	if result.Error !=nil{
		c.JSON(404, gin.H{
				"response": "Account doesn't exist",
			})
			return
	}
	err:=bcrypt.CompareHashAndPassword([]byte(findUser.Password),[]byte(user.Password))
	if err!=nil{
		c.JSON(401, gin.H{
				"response": "Invalid Password",
			})
			return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub":findUser.ID,
	"obj":findUser.Email,
	"exp": time.Now().Add(time.Hour*24*30).Unix(),
})

// Sign and get the complete encoded token as a string using the secret
tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNINING_KEY")))
if err!=nil{
	log.Fatal("failed to sign")
}
c.JSON(200,gin.H{
	"Token":tokenString,
})
}

func Validated(c *gin.Context){
	fmt.Println(c.Get("user"))
	c.JSON(200,gin.H{
	"response":"logged in",
  })
}