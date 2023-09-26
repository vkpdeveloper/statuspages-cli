package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestQuery struct {
	PerPage int
	pageId  string
}

func makeHttpRequest(method string, url string, body io.Reader, apiKey string, query RequestQuery, parsedResponse interface{}) error {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "OAuth "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()

	if query.PerPage > 0 {
		q.Add("per_page", fmt.Sprint(query.PerPage))
	}

	if query.pageId != "" {
		q.Add("page_id", query.pageId)
	}

	req.URL.RawQuery = q.Encode()

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &parsedResponse)

	if err != nil {
		return err
	}

	return nil
}

func HttpGet(url string, apiKey string, query RequestQuery, parsedResponse interface{}) error {
	return makeHttpRequest("GET", url, nil, apiKey, query, parsedResponse)
}
