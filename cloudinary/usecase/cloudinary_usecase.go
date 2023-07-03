package usecase

import (
	"warunk-bem/domain"
	"warunk-bem/helpers"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type media struct{}

func NewMediaUpload() domain.ClourdinaryUsecase {
	return &media{}
}

func (*media) FileUpload(file domain.File) (string, error) {
	//validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := helpers.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

func (*media) RemoteUpload(url domain.Url) (string, error) {
	//validate
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, errUrl := helpers.ImageUploadHelper(url.Url)
	if errUrl != nil {
		return "", err
	}
	return uploadUrl, nil
}
