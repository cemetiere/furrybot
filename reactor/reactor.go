package reactor

import (
	"errors"
	"fmt"
	"furrybot/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

const REACTOR_LINK = "https://furry.reactor.cc/random"

func scrapeForImageLink(r io.Reader) string {
	tokenizer := html.NewTokenizer(r)

	for {
		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			return ""
		} else if tt == html.StartTagToken {
			tagName, hasAttr := tokenizer.TagName()

			if tagName[0] != 'a' || !hasAttr {
				continue
			}

			// Parse attributes to see if this is
			// the link we're looking for
			// Extracts the link if true
			isLinkToImage := false
			linkToImage := ""
			for {
				key, val, moreAttr := tokenizer.TagAttr()

				if string(key) == "class" && string(val) == "prettyPhotoLink" {
					isLinkToImage = true
				}

				if string(key) == "href" {
					linkToImage = string(val)
				}

				if !moreAttr {
					break
				}
			}

			if isLinkToImage && len(linkToImage) > 0 {
				return linkToImage
			}
		}
	}
}

func fetchImage(folder, imageLink string) string {
	resp, err := http.Get(imageLink)

	if err != nil {
		log.Printf("Failed to download image. error: %s", err)
		return ""
	}
	defer resp.Body.Close()

	link_parts := strings.Split(imageLink, "/")
	filename := path.Join(folder, link_parts[len(link_parts)-1])

	out, err := os.Create(filename)

	if err != nil {
		log.Printf("failed to create file to store image. error: %s", err)
		return ""
	}
	defer out.Close()

	io.Copy(out, resp.Body)

	return filename
}

// Получает рандомную картинку в реакторе и сохраняет её в папке
func FetchFromReactor() (string, error) {
	resp, err := http.Get(REACTOR_LINK)

	if err != nil {
		return "", fmt.Errorf("failed to fetch from reactor. error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch from reactor. invalid status code: %v", resp.StatusCode)
	}

	link := scrapeForImageLink(resp.Body)

	if link == "" {
		return "", errors.New("couldn't find image link in HTML")
	}

	link = "http:" + link

	filename := fetchImage(config.Settings.ReactorFolderName, link)

	if filename == "" {
		return "", errors.New("image download failed")
	}

	return filename, nil
}
