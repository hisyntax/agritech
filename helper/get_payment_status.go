package helper

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/hisyntax/agritech/persistence"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hisyntax/monnify-go/transaction"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// get accepted payment status 	godoc
// @Summary      get accepted payment status
// @Description  this endpoint is used to comfirm the payment of the customers
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        ref  query  string  true  "ref"
// @Param        id   query  string  true  "id"
// @Success      200
// @Router       /payment/comfirm [get]
// func GetAcceptPaymentStatus(c *gin.Context) {
// ref := c.Query("ref")
// id := c.Query("id")
// if id == "" || ref == "" {
// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"error": persistence.ParamErr.Error(),
// 	})
// 	return
// }
// res, _, err := transaction.GetTransactionStatus(ref)
// if err != nil {
// 	fmt.Println(err)
// }

// if res.ResponseBody.PaymentStatus == "PAID" {
// 	// _id, _ := primitive.ObjectIDFromHex(id)
// 	// filter := bson.D{{Key: "_id", Value: _id}}

// } else if res.ResponseBody.PaymentStatus == "PENDING" {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": persistence.PaymentPending.Error(),
// 	})
// 	return
// } else {
// 	c.JSON(http.StatusOK, gin.H{
// 		"error": persistence.PaymentFailed.Error(),
// 	})
// 	return
// }
// }
