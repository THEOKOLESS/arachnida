package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"golang.org/x/net/html"
	"net/url"
)

// DownloadImages recursively downloads all images from the given URL into the destDir.
func DownloadImages(pageURL, destDir string, visited map[string]bool) error {
	if visited[pageURL] {
		return nil
	}
	visited[pageURL] = true

	resp, err := http.Get(pageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch %s: %s", pageURL, resp.Status)
	}

	base, err := url.Parse(pageURL)
	if err != nil {
		return err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) { 
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imgURL, err := base.Parse(attr.Val)
					if err == nil {
						downloadImage(imgURL.String(), destDir)
					}
				}
			}
		}
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					linkURL, err := base.Parse(attr.Val)
					if err == nil && strings.HasPrefix(linkURL.String(), base.Scheme+"://"+base.Host) {
						DownloadImages(linkURL.String(), destDir, visited)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nil
}

func downloadImage(imgURL, destDir string) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	_, file := path.Split(imgURL)
	if file == "" {
		return
	}
	os.MkdirAll(destDir, 0755)
	out, err := os.Create(path.Join(destDir, file))
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}
