package server_data_test

import "testing"
import "server_id_api/servers"

func TestGetServerDataIsOk(t *testing.T) {
  data := servers.GetServerData("google.com")

  if len(data.Endpoints) != 2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %d servers; should be 2", len(data.Endpoints))
  }

  server1 := data.Endpoints[0]
  server2 := data.Endpoints[1]

  expectedServer1 := servers.Server{
    Address: "172.217.5.110",
    SslGrade: "A",
    Country: "US",
    Owner:  "Google",
  }

  expectedServer2 := servers.Server{
    Address: "2607:f8b0:4005:808:0:0:0:200e",
    SslGrade: "A",
    Country: "US",
    Owner: "Google",
  }

  if server1 != expectedServer1 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %v server 1; should be %v", server1, expectedServer1)
  }

  if server2 != expectedServer2 {
    t.Errorf("TestGetServerDataIsOk(google.com) got %v server 2; should be %v", server2, expectedServer2)
  }
}
