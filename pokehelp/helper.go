package pokehelp

import (
	"io"
	"log"
	"net/http"
)

func GetBodyFromUrl(url string, _config *RequestConfig) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	// TODO: Handle 404 and other stuff
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code %d and body\n: %s\n", res.StatusCode, body)
	}

	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}

	return body, nil
}
