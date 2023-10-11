package images

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Getimages(context *gin.Context) {
	productId := context.Param("id")
	images, err := GetSpecificProductImage(productId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "Could not fetch products",
		})
	}
	context.JSON(http.StatusOK, gin.H{

		"success": true,
		"message": "Could not fetch products",
		"images":  images,
	})
}
func UploadMainimage(context *gin.Context, mainImageString string, productName string) (mainimagepath string, err error) {

	imageuuid := uuid.New()

	mainImageFilename := imageuuid.String() + productName

	imageBytes, err := base64.StdEncoding.DecodeString(mainImageString)

	if err != nil {
		return "", err
	}

	mainImagesFolder := "assets/products/"

	if _, err = os.Stat(mainImagesFolder); os.IsNotExist(err) {
		if err = os.Mkdir(mainImagesFolder, 0755); err != nil {
			return "", err
		}
	}
	productFolder := mainImagesFolder + strings.ReplaceAll(productName, " ", "")

	if _, err = os.Stat(productFolder); os.IsNotExist(err) {
		if err = os.Mkdir(productFolder, 0755); err != nil {
			return "", err
		}
	}

	mainImageFolder := productFolder + "/mainimage"

	if _, err = os.Stat(mainImageFolder); os.IsNotExist(err) {
		if err = os.Mkdir(mainImageFolder, 0755); err != nil {
			return "", err
		}
	}

	imagePath := filepath.Join(mainImageFolder, mainImageFilename)

	err = os.WriteFile(imagePath, imageBytes, 0644)
	if err != nil {
		return "", err
	}

	mainimagepath = imagePath

	return mainimagepath, nil
}
