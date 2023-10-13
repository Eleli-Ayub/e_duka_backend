package category

import (
	"errors"
	"net/http"

	globalutils "eleliafrika.com/backend/global_utils"
	"eleliafrika.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCategory(context *gin.Context) {
	var categoryInput models.Category

	if err := context.ShouldBindJSON(&categoryInput); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"category_input_error": err.Error(),
			"success":              false,
			"message":              "Wrong input from user",
		})
		return
	}
	success, err := ValidateCategoryInput(&categoryInput)
	if err != nil {
		response := models.Reply{
			Message: "error validating user input",
			Error:   err.Error(),
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else if !success {
		response := models.Reply{
			Message: "error validating user input for brand",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	} else {

		// check if category already exists
		category, err := FetchSingleCategory(categoryInput.CategoryName)
		if err != nil {
			response := models.Reply{
				Message: "error validating the request",
				Error:   err.Error(),
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else if category.IsDeleted {
			globalutils.HandleSuccess("category exists but is deleted. Reactivate the category", category, context)
			return
		} else if category.CategoryName != "" {

			response := models.Reply{
				Message: "category already exists",
				Data:    category,
				Success: false,
			}
			context.JSON(http.StatusBadRequest, response)
			return
		} else {
			categoryuuid := uuid.New()

			newProduct := models.Category{
				CategoryID:    categoryuuid.String(),
				CategoryName:  categoryInput.CategoryName,
				CategoryImage: categoryInput.CategoryImage,
			}

			category, err := newProduct.Save()

			if err != nil {

				globalutils.HandleError("coud not create category", err, context)
				return
			}

			globalutils.HandleSuccess("category created!!", category, context)
			return
		}

	}
}

func GetCategories(context *gin.Context) {
	categories, err := FetchAllCategories()
	if err != nil {
		response := models.Reply{
			Error:   err.Error(),
			Message: "error fetching categories",
			Success: false,
		}
		context.JSON(http.StatusBadRequest, response)
	} else {
		response := models.Reply{
			Message: "fetched categories succesful",
			Success: true,
			Data:    categories,
		}
		context.JSON(http.StatusOK, response)
	}
}

func DeleteCategory(context *gin.Context) {
	categoryname := context.Param("name")
	categoryExist, err := FetchSingleCategory(categoryname)
	if err != nil {
		globalutils.HandleError("error cheking validity of request", err, context)
		return
	} else if categoryExist.CategoryName == "" {
		globalutils.HandleError("the category requested is missing", errors.New("category doe not exist"), context)
		return
	} else if categoryExist.IsDeleted {
		globalutils.HandleError("the category requested is already deleted!!please confirm the validity of the request", errors.New("category already exist but is deleted"), context)
		return
	} else {
		deletedCategory, err := UpdateCategory(categoryname, models.Category{
			IsDeleted: true,
		})
		if err != nil {
			globalutils.HandleError("could not delete the product", err, context)
			return
		} else {
			globalutils.HandleSuccess("delete operation succesful!!", deletedCategory, context)
			return
		}
	}

}
