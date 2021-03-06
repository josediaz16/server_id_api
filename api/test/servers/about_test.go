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
    Owner:  "Google LLC",
  }

  expectedServer2 = model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Google LLC",
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

func getData(fixtureFile string, domainName string) (domain model.Domain) {
  server := mockRequest(fixtureFile)
  defer server.Close()

  apiClient := api.API{server.Client(), server.URL}
  return servers.GetServerData(&apiClient, domainName)
}

// Begin Test Cases

func TestGetServerDataServerDown(t *testing.T) {
  SetupTest()

  data := getData("../fixtures/ssllabs_down_server.json", "google.com")

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
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, ""},
    "ServersChanged":   []interface{}{data.ServersChanged, false},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  //Expectations
  runExpectations(t, expectations)
}

func TestGetServerDataServerOk(t *testing.T) {
  SetupTest()
  data := getData("../fixtures/ssllabs_server_ok.json", "google.com")

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
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, ""},
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

  data := getData("../fixtures/ssllabs_server_ok.json", "google.com")

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
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, "U"},
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

  data := getData("../fixtures/ssllabs_server_ok.json", "google.com")

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
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, "A"},
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
    Owner:  "Google LLC",
  }

  previousServer2 := model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Google LLC",
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

  data := getData("../fixtures/ssllabs_server_ok.json", "google.com")

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
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, "B"},
    "ServersChanged":   []interface{}{data.ServersChanged, true},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 2},
  }

  runExpectations(t, expectations)
}

func TestGetServerDataServerNotFound(t *testing.T) {
  SetupTest()

  data := getData("../fixtures/server_not_found.json", "dontexist.co")

  countDomains, _ := conn.Count("domains")
  countServers, _ := conn.Count("servers")

  var expectations = map[string][]interface{}{
    "servers":          []interface{}{len(data.Servers), 0},
    "Title":            []interface{}{data.Title, ""},
    "IsDown":           []interface{}{data.IsDown, true},
    "SslGrade":         []interface{}{data.SslGrade, "U"},
    "PreviousSslGrade": []interface{}{data.PreviousSslGrade, ""},
    "ServersChanged":   []interface{}{data.ServersChanged, false},
    "Domain Count":     []interface{}{countDomains, 1},
    "Server Count":     []interface{}{countServers, 0},
  }

  //Expectations
  runExpectations(t, expectations)
}

func TestGetAllDomainsNoDomainYet(t *testing.T) {
  SetupTest()
  domains, _ := servers.GetAllDomains()

  if len(domains) > 0 {
    t.Errorf("Domain count got %d, should be 0", len(domains))
  }
}

func TestGetAllDomainsDomainsExist(t *testing.T) {
  SetupTest()

  domainServer1 := model.Server{
    Address: "172.217.5.110",
    SslGrade: "B",
    Status: "Ready",
    Country: "US",
    Owner:  "Facebook",
  }

  domainServer2 := model.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Status: "Ready",
    Country: "US",
    Owner: "Facebook",
  }

  domain := model.Domain{
    Name: "facebook.com",
    SslGrade: "B",
    Title: "Facebook",
    Logo: "facebook.png",
    IsDown: true,
    Servers: []model.Server{domainServer1, domainServer2},
  }

  domain2 := model.Domain{
    Name: "google.com",
    SslGrade: "A",
    Title: "Google",
    Logo: "google.com.png",
    IsDown: false,
    Servers: []model.Server{expectedServer1, expectedServer2},
  }

  domain.InsertWithServers()

  domain2.InsertWithServers()

  domains, _ := servers.GetAllDomains()

  var expectations = map[string][]interface{}{
    "Domain Count":     []interface{}{len(domains), 2},
    "Title":            []interface{}{domains["google.com"].Title, "Google"},
    "Title2":           []interface{}{domains["facebook.com"].Title, "Facebook"},
    "ServerCount":      []interface{}{len(domains["google.com"].Servers), 2},
    "ServerCount2":     []interface{}{len(domains["facebook.com"].Servers), 2},
  }

  //Expectations
  runExpectations(t, expectations)
}


func TestWhoIsKnownServer(t *testing.T) {
  owner, country := servers.WhoIs("172.217.5.110")

  var expectations = map[string][]interface{}{
    "Owner":     []interface{}{owner, "Google LLC"},
    "Country":   []interface{}{country, "US"},
  }

  //Expectations
  runExpectations(t, expectations)
}

func TestWhoIsWeirdIp(t *testing.T) {
  owner, country := servers.WhoIs("2a03:2880:f127:283:face:b00c:0:25d")

  var expectations = map[string][]interface{}{
    "Owner":     []interface{}{owner, "Unknown"},
    "Country":   []interface{}{country, "Unknown"},
  }

  //Expectations
  runExpectations(t, expectations)
}
