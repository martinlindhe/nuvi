package news

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"gopkg.in/redis.v4"
)

const (
	cacheDir          = ".cache"
	redisNewsList     = "NEWS_XML"
	redisNewsHashList = "NEWS_XML_HASH"
)

var (
	redisClient *redis.Client
)

func init() {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err = os.Mkdir(cacheDir, 0777); err != nil {
			log.Fatal(err)
		}
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal("redis error:", err)
	}
}

// downloadAndOpenResource downloads resource if not found in local cache,
// then returns the data for the compressed files. Assumes url points to a zip
func downloadAndOpenResource(url string, verbose bool) (*zip.ReadCloser, error) {
	name := path.Base(url)
	cacheFile := path.Join(cacheDir, name)

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		if verbose {
			log.Println("saving", url, "to", cacheFile)
		}
		if err = downloadResource(url, cacheFile); err != nil {
			return nil, err
		}
	}

	r, err := zip.OpenReader(cacheFile)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func downloadResource(url, cacheFile string) error {
	out, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
