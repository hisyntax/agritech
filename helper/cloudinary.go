package helper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/joho/godotenv"
)

type Cloudinary struct {
	Cloud_Name string
	Api_Key    string
	Api_Secret string
	Folder     string
}

func NewCloudinary() *Cloudinary {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	name := os.Getenv("CLOUDINARY_CLOUD_NAME")
	api_key := os.Getenv("CLOUDINARY_API_KEY")
	api_secret := os.Getenv("CLOUDINARY_API_SECRET")
	folder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")

	return &Cloudinary{
		Cloud_Name: name,
		Api_Key:    api_key,
		Api_Secret: api_secret,
		Folder:     folder,
	}
}

func FileUploader(file interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cld := NewCloudinary()
	cldinary, err := cloudinary.NewFromParams(cld.Cloud_Name, cld.Api_Key, cld.Api_Secret)
	if err != nil {
		return "", err
	}

	//uploader
	uploaderP, err := cldinary.Upload.Upload(ctx, file, uploader.UploadParams{Folder: cld.Folder})
	if err != nil {
		return "", err
	}

	return uploaderP.SecureURL, nil
}
