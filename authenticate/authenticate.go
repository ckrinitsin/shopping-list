package authenticate

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ckrinitsin/shopping-list/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckAuth(c *gin.Context) {
	session := sessions.Default(c)
	token_session := session.Get("token")

	if token_session == nil {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		return
	}

	token_string, ok := token_session.(string)
	if !ok {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		return
	}

	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		c.Error(err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		return
	}

	var list models.List
	err = models.DB.
		Model(&models.List{}).
		Where("name = ?", claims["username"]).
		First(&list).
		Error
	if err != nil {
		c.Redirect(http.StatusFound, models.BasePath()+"/login")
		return
	}

	c.Set("current_list", list)

	c.Next()
}

func LoginGET(c *gin.Context) {
	title := "Shopping List"

	c.HTML(http.StatusOK, "login.html", gin.H{
		"name":      title,
		"error":     "",
		"base_path": models.BasePath(),
	})
}

func LoginPOST(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")

	var list models.List
	err := models.DB.
		Model(&models.List{}).
		Where("name = ?", username).
		First(&list).
		Error

	if err == gorm.ErrRecordNotFound {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "User does not exist",
		})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	err = bcrypt.CompareHashAndPassword(list.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Invalid username or password",
		})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	session := sessions.Default(c)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	}).SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	session.Set("token", token)
	session.Save()

	c.Redirect(http.StatusFound, models.BasePath()+"/")
}

func RegisterGET(c *gin.Context) {
	title := "Shopping List"

	c.HTML(http.StatusOK, "register.html", gin.H{
		"name":      title,
		"error":     "",
		"base_path": models.BasePath(),
	})
}

func RegisterPOST(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := c.PostForm("password")
	password_confirm := c.PostForm("password_confirm")
	global_password := strings.TrimSpace(c.PostForm("global_password"))

	if username == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Invalid username",
		})
		return
	}

	if len(password) <= 0 && len(password) <= 72 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Invalid password",
		})
		return
	}

	if password != password_confirm {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "The passwords do not match!",
		})
		return
	}

	if global_password != os.Getenv("GLOBAL_PASSWORD") {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "Global Password is wrong",
		})
		return
	}

	var count int64
	err := models.DB.
		Model(&models.List{}).
		Where("name = ?", username).
		Count(&count).
		Error

	if count > 0 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"error": "User does exist already",
		})
		return
	} else if err != gorm.ErrRecordNotFound && err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	var list models.List
	list = models.List{
		Name:     username,
		Password: hash,
	}

	err = models.DB.
		Create(&list).
		Error

	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"error": "Internal Server Error",
		})
		c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, models.BasePath()+"/login")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("token")
	session.Save()
	c.Redirect(http.StatusFound, models.BasePath()+"/login")
}
