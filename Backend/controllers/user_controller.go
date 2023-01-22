package controllers

import (
	"context"
	"net/http"
	"social/configs"
	"social/helpers"
	"social/models"
	"social/utils"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, utils.Responser{Status: http.StatusBadRequest, Message: err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, utils.Responser{Status: http.StatusBadRequest, Message: validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		password := helpers.HashPassword(*user.Password)
		newUser := models.User{
			ID:         primitive.NewObjectID(),
			First_name: user.First_name,
			Last_name:  user.Last_name,
			Email:      user.Email,
			Phone:      user.Phone,
			Password:   &password,
		}

		_, err = userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Responser{Status: http.StatusInternalServerError, Message: "error"})
			return
		}
		c.JSON(http.StatusCreated, utils.Responser{Status: http.StatusCreated, Message: "success", Data: newUser})
	}

}
