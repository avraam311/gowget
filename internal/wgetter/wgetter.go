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

func (wg *WGetter) WGet(link string, depth int) error {
	downloadedLinks := make(map[string]bool)

	link = strings.TrimRight(link, "/")
	URL, err := url.ParseRequestURI(link)
	if err != nil {
		return err
	}

	hostRe, err := regexp.Compile("https?://([a-z0-9]+[.])*\\.?(" + regexp.QuoteMeta(URL.Host) + "|iana.org)")
	if err != nil {
		return err
	}

	disallowedRe, err := regexp.Compile(`(?i)(mailto|tel|javascript):`)
	if err != nil {
		return err
	}

	if err := mkdir(downloadPath + "/" + URL.Host); err != nil {
		return err
	}

	collector := colly.NewCollector(
		colly.MaxDepth(depth),
		colly.Async(true),
		colly.URLFilters(hostRe),
		colly.DisallowedURLFilters(disallowedRe),
		colly.AllowedDomains(URL.Hostname()),
	)

	collector.OnHTML("a[href], link[rel=stylesheet], script[src], img[src], link[rel=icon]", func(el *colly.HTMLElement) {
		ul := el.Request.AbsoluteURL(el.Attr("href"))
		if !downloadedLinks[ul] {
			downloadedLinks[ul] = true
			if err := collector.Visit(ul); err != nil {
				log.Printf("visit %s failed: %v", ul, err)
			}
		}
	})

	collector.OnResponse(func(r *colly.Response) {
		u := r.Request.URL
		fullPath := downloadPath + "/" + URL.Hostname() + u.Path

		if u.RawQuery != "" {
			fullPath += "?" + u.RawQuery
		}

		fullPath = strings.TrimSuffix(fullPath, "/")

		if path.Ext(fullPath) == "" && !strings.HasSuffix(u.Path, "/") {
			fullPath += "/index.html"
		}

		if downloadedLinks[fullPath] {
			return
		}

		downloadedLinks[fullPath] = true

		if path.Ext(fullPath) == "" || strings.Contains(fullPath, "index.html") {
			dir := strings.TrimSuffix(fullPath, "index.html")
			if err := mkdir(dir); err != nil {
				log.Printf("mkdir %s failed: %v", dir, err)
			}
		} else {
			dirPath := path.Dir(fullPath)
			if err := mkdir(dirPath); err != nil {
				log.Printf("mkdir %s failed: %v", dirPath, err)
			}
		}

		if err = r.Save(fullPath); err != nil {
			log.Printf("save %s failed: %v", fullPath, err)
			return
		}

		fmt.Println("saved:", URL.Hostname()+u.Path)
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
