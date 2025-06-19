package utils

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{cld: cld}, nil
}

func (cs *CloudinaryService) UploadImage(file interface{}, folder string) (*uploader.UploadResult, error) {
	ctx := context.Background()
	
	uploadParams := uploader.UploadParams{
		Folder:           folder,
		UniqueFilename:   cloudinary.Bool(true),
		OverwriteExisting: cloudinary.Bool(false),
		ResourceType:     "image",
		Transformation:   "c_limit,w_1000,h_1000,q_auto:good",
	}

	result, err := cs.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *CloudinaryService) DeleteImage(publicID string) error {
	ctx := context.Background()
	
	_, err := cs.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	
	return err
}
