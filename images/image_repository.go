package images

import (
	"furrybot/reactor"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type IImageRepository interface {
	GetRandomImagePath() (string, error)
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

func (repository *LocalFilesImageRepository) GetRandomImagePath() (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(len(repository.availableImages))
	return path.Join(repository.imageFolder, repository.availableImages[randomIndex]), nil
}

func (repository *LocalFilesImageRepository) GetImages() []string {
	return repository.availableImages
}

type ReactorImageRepository struct{}

func (r *ReactorImageRepository) GetRandomImagePath() (string, error) {
	return reactor.FetchFromReactor()
}

func (r *ReactorImageRepository) GetImages() []string {
	return []string{}
}
