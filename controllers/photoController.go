package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"task-5-pbi-btpns-daffajatmiko/app"
	"task-5-pbi-btpns-daffajatmiko/app/auth"
	errorformat "task-5-pbi-btpns-daffajatmiko/helpers/error"
	"task-5-pbi-btpns-daffajatmiko/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetPhoto retrieves a list of photo profiles.
func GetPhoto(c *gin.Context) {
	// Create a list of photos
	photos := []models.Photo{}

	// Set the database
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&models.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	// Initialize the list of photos
	if len(photos) > 0 {
		for i := range photos {
			user := models.User{}
			err := db.Model(&models.User{}).Where("id = ?", photos[i].UserID).Take(&user).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": err.Error(),
					"data":    nil,
				})
				return
			}

			photos[i].Owner = app.Owner{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}
		}
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Data retrieved successfully",
		"data":    photos,
	})
}

// CreatePhoto creates a new photo profile.
func CreatePhoto(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Get the bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Token not found"})
		return
	}

	// Get the user email from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Get user data from the database
	var userHasLogin models.User

	err = db.Debug().Where("email = ?", email).First(&userHasLogin).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
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
	inputPhoto := models.Photo{}
	err = json.Unmarshal(body, &inputPhoto)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Initialize the photo
	inputPhoto.Init()
	inputPhoto.UserID = userHasLogin.ID
	inputPhoto.Owner = app.Owner{
		ID:       userHasLogin.ID,
		Username: userHasLogin.Username,
		Email:    userHasLogin.Email,
	}
	err = inputPhoto.Validate("upload") // Validate the photo
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Check if the photo already exists
	var oldPhoto models.Photo
	err = db.Debug().Model(&models.Photo{}).Where("user_id = ?", userHasLogin.ID).Find(&oldPhoto).Error
	if err != nil {
		if err.Error() == "Data not found" {
			err = db.Debug().Create(&inputPhoto).Error // Create the photo in the database
			if err != nil {
				formattedError := errorformat.ErrorMessage(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": formattedError.Error(),
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  "Success",
				"message": "Photo uploaded successfully",
				"data":    inputPhoto,
			})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Update the photo with new data
	inputPhoto.ID = oldPhoto.ID
	err = db.Debug().Model(&oldPhoto).Updates(&inputPhoto).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo changed successfully",
		"data":    inputPhoto,
	}) // Return the response
}

// UpdatePhoto updates a photo profile.
func UpdatePhoto(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Get the bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Token not found"})
		return
	}

	// Get the user email from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Get user data from JWT
	var userHasLogin models.User

	err = db.Debug().Where("email = ?", email).First(&userHasLogin).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
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
	photoInput := models.Photo{}
	err = json.Unmarshal(body, &photoInput)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Validate the photo
	err = photoInput.Validate("change")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Check if the photo already exists
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo with id " + c.Param("photoId") + " not found",
			"data":    nil,
		})
		return
	}

	// Validate the user ID
	if userHasLogin.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "You can't change the photo of another user",
			"data":    nil,
		})
		return
	}

	// Update the photo in the database
	err = db.Model(&photo).Updates(&photoInput).Error
	if err != nil {
		formattedError := errorformat.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": formattedError.Error(),
			"data":    nil,
		})
		return
	}

	photo.Owner = app.Owner{
		ID:       userHasLogin.ID,
		Username: userHasLogin.Username,
		Email:    userHasLogin.Email,
	}

	// Response for success
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo updated successfully",
		"data":    photo,
	})
}

// DeletePhoto deletes a photo profile.
func DeletePhoto(c *gin.Context) {
	// Set the database
	db := c.MustGet("db").(*gorm.DB)

	// Get the bearer token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "Token not found"})
		return
	}

	// Get the user email from JWT
	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Get user data from JWT
	var userHasLogin models.User
	if err := db.Debug().Where("email = ?", email).First(&userHasLogin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil,
		})
		return
	}

	// Check if the photo already exists
	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	// Validate the user ID
	if userHasLogin.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "You can't delete the photo of another user",
			"data":    nil,
		})
		return
	}

	// Delete the photo from the database
	err = db.Debug().Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo deleted successfully",
		"data":    nil,
	}) // Return the response
}
