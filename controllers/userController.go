package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"task-5-pbi-btpns-daffajatmiko/app"
	"task-5-pbi-btpns-daffajatmiko/app/auth"
	errorformat "task-5-pbi-btpns-daffajatmiko/helpers/error"
	"task-5-pbi-btpns-daffajatmiko/helpers/hash"
	"task-5-pbi-btpns-daffajatmiko/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login.
func Login(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Read the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Convert JSON to an object
	userModel := models.User{}
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Initialize user
	userModel.Init()
	err = userModel.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Check if the user exists
	var userLogin app.UserLogin
	err = db.Debug().Table("users").Select("*").Joins("LEFT JOIN photos ON photos.user_id = users.id").
		Where("users.email = ?", userModel.Email).Find(&userLogin).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + userModel.Email + " not found",
			"data":    nil,
		})
		return
	}

	// Verify the password
	err = hash.CheckPasswordHash(userLogin.Password, userModel.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	// Generate a token upon successful login
	token, err := auth.GenerateJWT(userLogin.Email, userLogin.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserData{
		ID:       userLogin.ID,
		Username: userLogin.Username,
		Email:    userLogin.Email,
		Token:    token,
		Photos: app.Photo{
			Title:    userLogin.Title,
			Caption:  userLogin.Caption,
			PhotoUrl: userLogin.PhotoUrl,
		},
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Login successfully",
		"data":    data,
	})
}

// CreateUser handles user registration.
func CreateUser(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Read the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Convert JSON to an object
	userModel := models.User{}
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	userModel.Init() // Initialize the user

	err = userModel.Validate("update") // Validate the user
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = userModel.HashPassword() // Hash the password
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Create(&userModel).Error // Create the user in the database
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserRegister{ // Data to be used for the response
		ID:        userModel.ID,
		Username:  userModel.Username,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User registered successfully",
		"data":    data,
	}) // Response for success
}

// UpdateUser handles user profile updates.
func UpdateUser(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Check if the user exists
	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + c.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Convert JSON to an object
	userModel := models.User{}

	userModel.ID = user.ID
	err = json.Unmarshal(body, &userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Validate the user
	err = userModel.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Hash the password
	err = userModel.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	// Update the user
	err = db.Debug().Model(&user).Updates(&userModel).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserRegister{ // Data to be used for the response
		ID:        userModel.ID,
		Username:  userModel.Username,
		Email:     userModel.Email,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}

	// Response for success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Error",
		"message": "User updated successfully",
		"data":    data,
	})
}

// DeleteUser handles user deletion.
func DeleteUser(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Check if the user exists
	var user models.User

	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + c.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	// Delete the user
	err = db.Debug().Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Response for success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User deleted successfully",
		"data":    nil,
	})
}
