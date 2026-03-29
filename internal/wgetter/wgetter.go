package wgetter

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const (
	downloadPath = "downloadedSite"
)

type WGetter struct {
}

func New() *WGetter {
	return &WGetter{}
}

func (wg *WGetter) WGet(link string) error {
	downloadedLinks := make(map[string]bool)

	link = strings.TrimRight(link, "/")
	URL, err := url.ParseRequestURI(link)
	if err != nil {
		return err
	}

	re, err := regexp.Compile("https?://([a-z0-9]+[.])*" + URL.Host)
	if err != nil {
		return err
	}

	if err := mkdir(downloadPath + "/" + URL.Host); err != nil {
		return err
	}

	collector := colly.NewCollector(colly.MaxDepth(1), colly.URLFilters(re))

	collector.OnHTML("a[href]", func(el *colly.HTMLElement) {
		ul := el.Request.AbsoluteURL(el.Attr("href"))
		if !downloadedLinks[ul] {
			if err := collector.Visit(ul); err != nil {
				return
			}
		}
	})

	collector.OnResponse(func(r *colly.Response) {
		reqURLPath := r.Request.URL.Path
		fullPath := downloadPath + "/" + URL.Hostname() + reqURLPath

		if downloadedLinks[fullPath] {
			return
		}

		downloadedLinks[fullPath] = true
		if path.Ext(fullPath) == "" {
			if err := mkdir(fullPath); err != nil {
				log.Printf("mkdir %s failed: %v", fullPath, err)
			}
		} else {
			dirPath := fullPath[:strings.LastIndexByte(fullPath, '/')]
			if err := mkdir(dirPath); err != nil {
				log.Printf("mkdir %s failed: %v", dirPath, err)
			}
		}

		if path.Ext(reqURLPath) == "" {
			if fullPath[len(fullPath)-1] != '/' {
				fullPath += "/"
			}
			fullPath += "index.html"
			if _, err := os.Create(fullPath); err != nil {
				fmt.Printf("error creating file: %s\n", err.Error())
			}
		}

		if err = r.Save(fullPath); err != nil {
			log.Printf("save %s failed: %v", fullPath, err)
			return
		}

		fmt.Println("saved:", URL.Hostname()+reqURLPath)
	})

	if err = collector.Visit(URL.String()); err != nil {
		return fmt.Errorf("visit %s: %w", URL.String(), err)
	}
	collector.Wait()
	return nil
}

func mkdir(folderName string) error {
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir %s: %w", folderName, err)
	}
	return nil
}
