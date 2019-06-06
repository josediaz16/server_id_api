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
  "time"
)
var conn *db.Conn = db.NewConn()

var expectedServer1, expectedServer2 model.Server

func SetupTest() {
  conn.DeleteAllFrom("domains")
  expectedServer1 = model.Server{
    Address: "172.217.5.110",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner:  "Google",
  }

  expectedServer2 = model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Google",
  }
}

func runExpectations(t *testing.T, expectations map[string][]interface{}) {
  for key, pair := range expectations {

    if pair[0] != pair[1] {
      t.Errorf("TestGetServerDataIsOk(google.com) got %s %v; should be %v", key, pair[0], pair[1])
    }
  }
}

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

func getData(fixtureFile string) (domain model.Domain) {
  server := mockRequest(fixtureFile)
  defer server.Close()

  apiClient := api.API{server.Client(), server.URL}
  return servers.GetServerData(&apiClient, "google.com")
}

// Begin Test Cases

func TestGetServerDataServerDown(t *testing.T) {
  SetupTest()

  data := getData("../fixtures/ssllabs_down_server.json")

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  expectedServer2.SslGrade = "U"
  expectedServer2.Status = "Unable to connect to the server"

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  var expectations = map[string][]interface{}{
    "server 1":         []interface{}{server1, expectedServer1},
    "server 2":         []interface{}{server2, expectedServer2},
    "Title":            []interface{}{data.Title, "Google"},
    "IsDown":           []interface{}{data.IsDown, true},
    "SslGrade":         []interface{}{data.SslGrade, "U"},
    "ServersChanged":   []interface{}{data.ServersChanged, false},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  //Expectations
  runExpectations(t, expectations)
}

func TestGetServerDataServerOk(t *testing.T) {
  SetupTest()
  data := getData("../fixtures/ssllabs_server_ok.json")

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  var expectations = map[string][]interface{}{
    "server 1":         []interface{}{server1, expectedServer1},
    "server 2":         []interface{}{server2, expectedServer2},
    "Title":            []interface{}{data.Title, "Google"},
    "IsDown":           []interface{}{data.IsDown, false},
    "SslGrade":         []interface{}{data.SslGrade, "A"},
    "ServersChanged":   []interface{}{data.ServersChanged, false},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  //Expectations
  runExpectations(t, expectations)
}

func TestGetServerDataServerChanged(t *testing.T) {
  SetupTest()
  domain := model.Domain{
    Name: "google.com",
    SslGrade: "U",
    Title: "Google",
    Logo: "mylogo.png",
    IsDown: true,
  }

  domain.Insert()
  domain.UpdatedAt = time.Now().Add(-2*time.Hour).Format(model.TimeLayout)
  domain.Update()

  data := getData("../fixtures/ssllabs_server_ok.json")

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  //Expectations

  var expectations = map[string][]interface{}{
    "server 1":         []interface{}{server1, expectedServer1},
    "server 2":         []interface{}{server2, expectedServer2},
    "Title":            []interface{}{data.Title, "Google"},
    "IsDown":           []interface{}{data.IsDown, false},
    "SslGrade":         []interface{}{data.SslGrade, "A"},
    "ServersChanged":   []interface{}{data.ServersChanged, true},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  runExpectations(t, expectations)
}

func TestGetServerDataServerNotChanged(t *testing.T) {
  SetupTest()
  domain := model.Domain{
    Name: "google.com",
    SslGrade: "A",
    Title: "Google",
    Logo: "mylogo.png",
    IsDown: true,
  }

  domain.Insert()
  domain.UpdatedAt = time.Now().Add(-2*time.Hour).Format(model.TimeLayout)
  domain.Update()

  data := getData("../fixtures/ssllabs_server_ok.json")

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  //Expectations

  var expectations = map[string][]interface{}{
    "server 1":         []interface{}{server1, expectedServer1},
    "server 2":         []interface{}{server2, expectedServer2},
    "Title":            []interface{}{data.Title, "Google"},
    "IsDown":           []interface{}{data.IsDown, false},
    "SslGrade":         []interface{}{data.SslGrade, "A"},
    "ServersChanged":   []interface{}{data.ServersChanged, false},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  runExpectations(t, expectations)
}

func TestGetServerDataBetweenTime(t *testing.T) {
  SetupTest()

  previousServer1 := model.Server{
    Address: "172.217.5.110",
    SslGrade: "B",
    Status: "Ready",
    Country: "US",
    Owner:  "Google",
  }

  previousServer2 := model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Google",
  }

  domain := model.Domain{
    Name: "google.com",
    SslGrade: "B",
    Title: "Google",
    Logo: "mylogo.png",
    IsDown: true,
    Servers: []model.Server{previousServer1, previousServer2},
  }

  domain.InsertWithServers()
  domain.UpdatedAt = time.Now().Add(-58*time.Minute).Format(model.TimeLayout)
  domain.Update()

  data := getData("../fixtures/ssllabs_server_ok.json")

  server1 := data.Servers[0]
  server2 := data.Servers[1]

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  //Expectations

  var expectations = map[string][]interface{}{
    "server 1":         []interface{}{server1, expectedServer1},
    "server 2":         []interface{}{server2, expectedServer2},
    "Title":            []interface{}{data.Title, "Google"},
    "IsDown":           []interface{}{data.IsDown, false},
    "SslGrade":         []interface{}{data.SslGrade, "A"},
    "ServersChanged":   []interface{}{data.ServersChanged, true},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  runExpectations(t, expectations)
}
