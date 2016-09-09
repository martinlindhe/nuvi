package news

import (
	"bytes"
	"io"
	"log"
)

// ImportPublishedDocuments imports xml from all published zips to redis
func ImportPublishedDocuments(verbose bool) error {
	listing, err := fetchFileListing()
	if err != nil {
		return err
	}

	for _, file := range listing {
		if err := importResource(file, verbose); err != nil {
			log.Println("error: failed to import", file, ":", err)
		}
	}
	return nil
}

func importResource(url string, verbose bool) error {
	r, err := downloadAndOpenResource(url, verbose)
	if err != nil {
		return err
	}
	defer r.Close()

	log.Println("importing", url)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		var b bytes.Buffer
		if _, err = io.Copy(&b, rc); err != nil {
			return err
		}
		rc.Close()

		// attach all xml content to redis list,
		// keep track of published content using a hash list
		_, err = redisClient.HGet(redisNewsHashList, f.Name).Result()
		if err == nil {
			if verbose {
				log.Println("redis: skipping imported", f.Name)
			}
			continue
		}

		if verbose {
			log.Println("redis: importing", f.Name)
		}
		redisClient.HSet(redisNewsHashList, f.Name, "1")
		redisClient.RPush(redisNewsList, b.String())
	}
	return nil
}
