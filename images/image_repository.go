package images

import (
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type IImageRepository interface {
	GetRandomImagePath() string
	GetImages() []string
}

type LocalFilesImageRepository struct {
	availableImages []string
	imageFolder     string
}

func NewLocalFilesImageRepository(imageFolder string) (*LocalFilesImageRepository, error) {
	repository := LocalFilesImageRepository{}

	files, err := os.ReadDir(imageFolder)

	if err != nil {
		return &repository, err
	}

	for _, v := range files {
		if v.Type().IsRegular() && (strings.HasSuffix(v.Name(), ".jpg") || strings.HasSuffix(v.Name(), ".png")) {
			repository.availableImages = append(repository.availableImages, v.Name())
		}
	}

	repository.imageFolder = imageFolder

	return &repository, nil
}

func (repository *LocalFilesImageRepository) GetRandomImagePath() string {
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(len(repository.availableImages))
	return path.Join(repository.imageFolder, repository.availableImages[randomIndex])
}

func (repository *LocalFilesImageRepository) GetImages() []string {
	return repository.availableImages
}
