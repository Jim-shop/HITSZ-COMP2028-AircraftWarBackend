package modules

import (
	"imshit/aircraftwar/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	User     string `form:"user" json:"user" uri:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

// 核验身份并签发token
func Login(c *gin.Context) {
	request := LoginRequest{}
	if err := c.Bind(&request); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user, err := models.QueryUser(request.User)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	token := models.NewToken(user)
	c.String(http.StatusOK, token.Token)
}