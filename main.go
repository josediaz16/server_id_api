package main

import (
  "server_id_api/servers"
  "server_id_api/api"
  "net/http"
)

func main() {
  var sslLabsClient = api.API{&http.Client{}, "https://api.ssllabs.com"}
  servers.GetServerData(&sslLabsClient, "google.com")
}
