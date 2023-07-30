package images

import (
	"errors"
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
	return reactor.FetchFromReactorRandom(reactor.REACTOR_LINK)
}

func (r *ReactorImageRepository) GetImages() []string {
	return []string{}
}

type FapReactorImageRepository struct{}

func (r *FapReactorImageRepository) GetRandomImagePath() (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	pageNo := rand.Intn(100) + 1
	links, err := reactor.ScrapeReactorTagPage(reactor.FAP_REACTOR_LINK, "Йифф", pageNo)

	if err != nil {
		return "", err
	}

	if len(links) == 0 {
		return "", errors.New("failed to find any images on search page")
	}

	postNo := rand.Intn(len(links))

	filename, err := reactor.FetchImageFromReactorPost(links[postNo])

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (r *FapReactorImageRepository) GetImages() []string {
	return []string{}
}
