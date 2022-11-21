package user

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hisyntax/agritech/helper"
	"github.com/hisyntax/agritech/persistence"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"

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

type Phone struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

// request otp 	godoc
// @Summary      request otp
// @Description  this endpoint is used to request an otp for a user
// @Tags         user
// @Accept       json
// @Produce      json
// @param        user  body  Phone  true  "user"
// @Success      200
// @Router       /user/otp [post]
func (u User) SendOTP(c *gin.Context) {
	var phone Phone

	if err := c.BindJSON(&phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error binding json",
		})
		return
	}
	if len(phone.PhoneNumber) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the first 0 must be taken off the number",
		})
		return
	}

	filter := bson.M{"email": phone.Email}
	_, userErr := persistence.GetMongoDoc(persistence.UserCollection, filter)
	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	data := bson.D{{Key: "phone_number", Value: phone.PhoneNumber}}

	updateRes, updateErr := persistence.UpdateMongoDoc(persistence.UserCollection, filter, data)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error adding phone number",
		})
		return
	}

	//generate otp
	otp, _ := helper.GenerateRandomGenerator(4)
	strOtp := fmt.Sprintf("%v", otp)
	if err := persistence.SetRedisValue(phone.Email, strOtp, 5*time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error geerating otp",
		})
		return
	}

	newNumber := "+234" + phone.PhoneNumber
	//set twillo the variavles
	accountSid := os.Getenv("TWILLO_ACCT_SID")
	authToken := os.Getenv("TWILLO_AUTH_KEY")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	msg := fmt.Sprintf("this is your otp %v. it expires in 5 minutes ", otp)
	params.SetBody(msg)
	fromNumber := os.Getenv("TWILLO_NUMBER")
	fmt.Printf("This is twillo number %v\n", fromNumber)
	params.SetFrom(fromNumber)
	params.SetTo(newNumber)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error sending otp",
		})
		return
	}

	response := fmt.Sprintf("otp has been sent to %v", newNumber)
	c.JSON(http.StatusOK, gin.H{
		"message":  response,
		"response": resp,
		"update":   updateRes,
	})
}

type OtpVal struct {
	Otp   int    `json:"otp"`
	Email string `json:"email"`
}

// validate otp 	godoc
// @Summary      validate otp
// @Description  this endpoint is used to validate an otp for a user
// @Tags         user
// @Accept       json
// @Produce      json
// @param        user  body  OtpVal  true  "user"
// @Success      200
// @Router       /user/otp/validate [post]
func (u User) ValidateOtp(c *gin.Context) {
	var otp OtpVal

	if err := c.BindJSON(&otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error binding json",
		})
		return
	}

	val, err := persistence.GetRedisValue(otp.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "otp has expired",
		})
		return
	}

	ott := fmt.Sprintf("%v", otp.Otp)
	if val != ott {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid otp",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "otp confirmed",
	})
}
