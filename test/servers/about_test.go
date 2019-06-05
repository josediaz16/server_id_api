package servers_test

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

func mockRequest(fixtureFile string) (*httptest.Server) {
  mockedResponse := getServerResponse(fixtureFile)

  return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    rw.Write([]byte(mockedResponse))
  }))
}

func TestGetServerDataServerDown(t *testing.T) {
  server := mockRequest("../fixtures/ssllabs_down_server.json")
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
    SslGrade: "U",
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

func TestGetServerDataServerOk(t *testing.T) {
  server := mockRequest("../fixtures/ssllabs_server_ok.json")
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
    SslGrade: "A",
    Status: "Ready",
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

  if data.IsDown {
    t.Errorf("TestGetServerDataIsOk(google.com) got IsDown true; should be false")
  }
}
