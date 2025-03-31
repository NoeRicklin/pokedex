package utils

import(
	"net/http"
	"encoding/json"
)

type AreaJSON struct {
	Count		int
	Next		string
	Previous	string
	Results		[]struct{
		Name	string
		Url		string
	}
}

func GetURLBody(url string) (AreaJSON, error) {
	res, err := http.Get(url)
	if err != nil {
		return AreaJSON{}, err
	}
	defer res.Body.Close()

	var body AreaJSON
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return AreaJSON{}, err
	}

	return body, nil
}

