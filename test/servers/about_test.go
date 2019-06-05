package servers_test

import (
  "testing"
  "io/ioutil"
  "os"
  "server_id_api/servers"
  "server_id_api/model"
  "server_id_api/api"
  "server_id_api/db"
  "net/http"
  "net/http/httptest"
)
var conn *db.Conn = db.NewConn()

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

func SetupTest() {
  conn.DeleteAllFrom("domains")
}

func TestGetServerDataServerDown(t *testing.T) {
  SetupTest()
  server := mockRequest("../fixtures/ssllabs_down_server.json")
  defer server.Close()

  apiClient := api.API{server.Client(), server.URL}
  data := servers.GetServerData(&apiClient, "google.com")

  if len(data.Servers) != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers; should be 2", len(data.Servers))
  }

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  expectedServer1 := model.Server{
    Address: "172.217.5.110",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner:  "Google",
  }

  expectedServer2 := model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "U",
    Status: "Unable to connect to the server",
    Country: "US",
    Owner: "Google",
  }

  //Expectations

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

  if data.SslGrade != "U" {
    t.Errorf("TestGetServerDataIsOk(google.com) got SslGrade %s; should be U", data.SslGrade)
  }

  countDomains, _ := conn.Count("domains")

  if countDomains != 1 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d domains in DB; should be 1", countDomains)
  }

  countServers, _ := conn.Count("servers")

  if countServers != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers in DB; should be 2", countServers)
  }
}

func TestGetServerDataServerOk(t *testing.T) {
  SetupTest()
  server := mockRequest("../fixtures/ssllabs_server_ok.json")
  defer server.Close()

  apiClient := api.API{server.Client(), server.URL}
  data := servers.GetServerData(&apiClient, "google.com")

  if len(data.Servers) != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers; should be 2", len(data.Servers))
  }

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  expectedServer1 := model.Server{
    Address: "172.217.5.110",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner:  "Google",
  }

  expectedServer2 := model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Google",
  }

  //Expectations

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

  if data.SslGrade != "A" {
    t.Errorf("TestGetServerDataIsOk(google.com) got SslGrade %s; should be A", data.SslGrade)
  }

  countDomains, _ := conn.Count("domains")

  if countDomains != 1 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d domains in DB; should be 1", countDomains)
  }

  countServers, _ := conn.Count("servers")

  if countServers != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers in DB; should be 2", countServers)
  }
}
