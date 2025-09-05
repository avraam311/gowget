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

	mkdir(downloadPath + "/" + URL.Host)

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
			mkdir(fullPath)
		} else {
			mkdir(fullPath[:strings.LastIndexByte(fullPath, '/')])
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
			panic(err)
		}

		fmt.Println("saved:", URL.Hostname()+reqURLPath)
	})

	if err = collector.Visit(URL.String()); err != nil {
		log.Panic("err: visit: " + err.Error())
	}
	collector.Wait()
	return nil
}

func mkdir(folderName string) {
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) && os.MkdirAll(folderName, os.ModePerm) != nil {
		log.Panic(err)
	}
}
