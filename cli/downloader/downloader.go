package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/alde/horus/cli/config"
)

// Downloader struct
type Downloader struct {
	c   *http.Client
	cfg *config.Config
}

// New creates a new Downloader
func New(client *http.Client, conf *config.Config) *Downloader {
	return &Downloader{
		c:   client,
		cfg: conf,
	}
}

// Download will download the secret
func (d *Downloader) Download(repo, key string) string {
	repo1 := url.QueryEscape(repo)
	r, err := d.c.Get(fmt.Sprintf("%s/api/v1/secret?repo=%s&key=%s", d.cfg.Horus.Host, repo1, key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(body)
}
