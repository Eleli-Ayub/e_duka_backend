package images

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"eleliafrika.com/backend/database"
	"eleliafrika.com/backend/models"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func GetProductImages() ([]models.ProductImage, error) {
	var productsImages []models.ProductImage
	err := database.Database.Find(&productsImages).Error
	if err != nil {
		return []models.ProductImage{}, err

	}
	return productsImages, nil
}

func GetSpecificProductImage(productid string) ([]string, error) {
	var productsImages []models.ProductImage
	var images []string
	err := database.Database.Where("product_id=?", productid).Find(&productsImages).Error
	if err != nil {
		return []string{}, err
	}
	for image := range productsImages {
		images = append(images, productsImages[image].ImageUrl)
	}
	return images, nil
}
func UploadImageToFireBase(imageString string) (string, error) {
	config := &firebase.Config{
		StorageBucket: "eduka-f19e5.appspot.com",
	}

	opt := option.WithCredentialsFile("./key.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}

	decodedImage, err := base64.StdEncoding.DecodeString(imageString)

	if err != nil {
		return "", err
	}

	// Generate a unique image name or use an existing one
	imageName := "newimage.jpg"

	object := bucket.Object("images/" + imageName)
	wc := object.NewWriter(context.Background())
	wc.ContentType = "image/jpeg"

	_, err = wc.Write(decodedImage)
	if err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	// Get the download URL
	downloadURL, err := object.Attrs(context.Background())
	if err != nil {
		return "", err
	}

	return downloadURL.MediaLink, nil

}

func UploadHandler(productName string, imageString string, context *gin.Context) (string, error) {
	imageData, err := base64.StdEncoding.DecodeString(imageString)

	if err != nil {
		return "", err
	}

	imageUUID := uuid.New()
	err = Uploader.UploadFile(bytes.NewReader(imageData), strings.ReplaceAll(productName, " ", "")+"/"+imageUUID.String())
	if err != nil {
		return "", err
	}

	imageUrl := fmt.Sprintf("https://storage.googleapis.com/%s/eduka/product_images/%s",
		BucketName,
		strings.ReplaceAll(productName, " ", "")+"/"+imageUUID.String())

	return imageUrl, nil
}

const (
	projectID  = "eduka-404606" // FILL IN WITH YOURS
	BucketName = "eduka-bucket" // FILL IN WITH YOURS
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var Uploader *ClientUploader

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/eleli/.config/gcloud/application_default_credentials.json") // FILL IN WITH YOUR FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	Uploader = &ClientUploader{
		cl:         client,
		bucketName: BucketName,
		projectID:  projectID,
		uploadPath: "eduka/product_images/",
	}

}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
