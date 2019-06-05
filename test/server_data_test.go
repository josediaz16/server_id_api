package server_data_test

import (
  "testing"
  "io/ioutil"
  "os"
  "server_id_api/servers"
  "server_id_api/api"
  "net/http"
  "net/http/httptest"
)

func getServerResponse(fixtureFile string) ([]byte) {
  jsonFile, _ := os.Open(fixtureFile)
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)
  return byteValue
}

func TestGetServerDataIsOk(t *testing.T) {
  mockedResponse := getServerResponse("fixtures/ssllabs_down_server.json")

  server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    rw.Write([]byte(mockedResponse))
  }))

  defer server.Close()

  apiClient := api.API{server.Client(), server.URL}

  data := servers.GetServerData(&apiClient, "google.com")

  if len(data.Servers) != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers; should be 2", len(data.Servers))
  }

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  expectedServer1 := servers.Server{
    Address: "172.217.5.110",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner:  "Google",
  }

  expectedServer2 := servers.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "",
    Status: "Unable to connect to the server",
    Country: "US",
    Owner: "Google",
  }

  if server1 != expectedServer1 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %v server 1; should be %v", server1, expectedServer1)
  }

  if server2 != expectedServer2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %v server 2; should be %v", server2, expectedServer2)
  }

  if data.Title != "Google" {
    t.Errorf("TestGetServerDataIsOk(google.com) got title %v; should be Google", data.Title)
  }

  if !data.IsDown {
    t.Errorf("TestGetServerDataIsOk(google.com) got IsDown false; should be true")
  }
}

func TestGetDomainHeadIsOk(t *testing.T) {
  title, logo := servers.GetDomainHead("stackoverflow.com")

  expectedTitle := "Stack Overflow - Where Developers Learn, Share, & Build Careers"
  expectedLogo := "https://cdn.sstatic.net/Sites/stackoverflow/img/favicon.ico?v=4f32ecc8f43d"

  if title != expectedTitle {
    t.Errorf("TestGetDomainHeadIsOk(stackoverflow.com) got title %v, should be %v", title, expectedTitle)
  }

  if logo != expectedLogo {
    t.Errorf("TestGetDomainHeadIsOk(stackoverflow.com) got logo %v, should be %v", logo, expectedLogo)
  }
}
