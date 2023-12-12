package cdn

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/joho/godotenv"
)

func Base64toEncode(bytes []byte) (string,error){
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += ToBase64(bytes)

	// Print the full base64 representation of the image
	fmt.Println(base64Encoding)
	return base64Encoding,nil
}

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}


func ImageKit(ctx context.Context, base64Image string) (string, error) {
	err := godotenv.Load("../.env")
    if err != nil {
        panic("Failed to load .env file")
    }
	
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
    publicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
    urlEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	fmt.Println("start uploading image ...")
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	
	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: urlEndpoint,
	})

	resp, err := ik.Uploader.Upload(ctx, base64Image, uploader.UploadParam{
		FileName: "test_image.jpg",
		Tags:     "annisa-picture",
		Folder:   "annisa-salon",
	})

	if err != nil {
		fmt.Printf("an error occurred when uploading image %v", err)
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("an error occurred when uploading image %v", resp)
		return "", errors.New("failed to upload image")
	}

	// Return the ImageKit URL
	return resp.Data.Url, nil
}