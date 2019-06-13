package api

import (
  "io/ioutil"
  "net/http"
)

type API struct {
  Client  *http.Client
  BaseURL string
}

func (api *API) GetWithParams(path string, params map[string]string) ([]byte, error) {
  req, _ := http.NewRequest("GET", api.BaseURL + path, nil)
  query := req.URL.Query()

  for key, value := range params {
    query.Add(key, value)
  }

  req.URL.RawQuery = query.Encode()
  resp, err := api.Client.Do(req)

  if err != nil {
    return nil, err
  }

  body, err := ioutil.ReadAll(resp.Body)

  return body, err
}
