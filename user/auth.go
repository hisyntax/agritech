package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/hisyntax/agritech/helper"
	"github.com/hisyntax/agritech/persistence"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SwaggerUserSignup struct {
	Full_Name string `json:"full_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// create user 	godoc
// @Summary      create user
// @Description  this endpoint is used create a user
// @Tags         user
// @Accept       json
// @Produce      json
// @param        user  body  SwaggerUserSignup  true  "user"
// @Success      200
// @Router       /user/signup [post]
func (user User) Signup(c *gin.Context) {
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := persistence.Validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//check if the email address conains @ and .
	if strings.Contains(user.Email, "@") && strings.Contains(user.Email, ".") {

		emailFilter := bson.D{{Key: "email", Value: user.Email}}
		_, emailErr := persistence.GetMongoDoc(persistence.UserCollection, emailFilter)
		if emailErr != nil {
			if err := user.BeforeCreate(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			userEmail := strings.ToLower(user.Email)
			user.Email = userEmail
			user.Created_At = time.Now()
			user.Updated_At = time.Now()
			user.ID = primitive.NewObjectID()

			_, insertErr := persistence.CreateMongoDoc(persistence.UserCollection, user)
			if insertErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": insertErr.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user": user.PublicUser(&user),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": persistence.ErrEmailTaken.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email address is not valid",
		})
		return
	}

}

type SwaggerUserSignin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// signin user 	godoc
// @Summary      signin user
// @Description  this endpoint is used signin a user
// @Tags         user
// @Accept       json
// @Produce      json
// @param        user  body  SwaggerUserSignin  true  "user"
// @Success      200
// @Router       /user/signin [post]
func (user User) Signin(c *gin.Context) {
	var u UserSignin
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := persistence.Validate.Struct(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userEmail := strings.ToLower(u.Email)
	u.Email = userEmail

	emailFilter := bson.D{{Key: "email", Value: u.Email}}
	foundUser, emailErr := user.GetUserDoc(persistence.UserCollection, emailFilter)
	if emailErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": persistence.ErrEmailNotFound.Error(),
		})
		return
	} else {
		if err := helper.VerifyPassword(foundUser.Password, u.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": persistence.ErrPasswordIncorrect.Error(),
			})
			return
		}
	}

	//generate a token for the user on signup
	token, _, _ := persistence.GenerateAllTokens(foundUser.ID.Hex())
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.PublicUser(foundUser),
	})
}
