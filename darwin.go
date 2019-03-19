package BingWallpaperGo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type wp struct {
	url         string
	path        string
	title       string
	description string
}

func SetBingWp(host, path string) (err error) {
	w := &wp{}
	if err = w.parse(host); err != nil {
		return err
	}
	if err = w.save(path); err != nil {
		return err
	}
	if err = w.set(); err != nil {
		return err
	}
	return nil
}

func (w *wp) parse(host string) (err error) {
	u := &url.URL{
		Host:   host,
		Scheme: "http",
	}
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	link, _ := doc.Find("[href*='1920x1080.jpg']").Attr("href")

	link = strings.Split(link, "&")[0]

	u, _ = url.Parse(link)
	u.Host = host
	u.Scheme = "http"

	u.Host = host
	w.url = u.String()
	return nil
}

func (w *wp) save(dir string) (err error) {
	res, err := http.Get(w.url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	fileName := ""
	if len(strings.Split(w.url, "=")) > 1 {
		fileName = strings.Split(w.url, "=")[1]
	}
	filePath := path.Join(dir, time.Now().Format("2006-01-02-")+fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	w.path = filePath
	return nil
}

func (w *wp) set() (err error) {
	err = exec.Command("osascript",
		"-e", `tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(w.path)).Run()
	if err != nil {
		return err
	}
	text := "text"

	err = exec.Command("osascript",
		"-e", fmt.Sprintf(`display notification "%s" with title "%s"`, text, w.title)).Run()
	return err
}
