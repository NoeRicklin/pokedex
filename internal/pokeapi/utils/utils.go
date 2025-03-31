package utils

import(
	"net/http"
	"time"
	"github.com/NoeRicklin/pokedex/internal/pokecache"
	"encoding/json"
	"io"
)

var c *pokecache.Cache
func SetupCache(dur time.Duration) {
	c = pokecache.NewCache(dur)
}

func GetURLBody[T any](url string) (T, error) {
	var rawBody []byte

	_, cacheHit := c.Get(url)
	if cacheHit {
		rawBody, _ = c.Get(url)
	} else {
		res, err := http.Get(url)
		if err != nil {
			var nilT T
			return nilT, err
		}
		defer res.Body.Close()

		rawBody, err = io.ReadAll(res.Body)
		if err != nil {
			var nilT T
			return nilT, err
		}
	}

	var output T
	if err := json.Unmarshal(rawBody, &output); err != nil {
			var nilT T
			return nilT, err
	}

	return output, nil
}

