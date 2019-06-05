package servers

import (
  "fmt"
  "log"
  "encoding/json"
  "os/exec"
  "strings"
  "server_id_api/api"
)

type Server struct {
  Address  string `json:"ipAddress"`
  SslGrade string `json:"grade"`
  Status   string `json:"statusMessage"`
  Country  string
  Owner    string
}

type SslLabsResponse struct {
  Servers   []Server  `json:"endpoints"`
  Title     string
  Logo      string
  IsDown    bool
}

const WhoIsCmd = "whois %s | grep -E \\(Country\\|OrgName\\) | awk '{print $2}' | xargs"

func GetServerData(apiClient *api.API, domain string) SslLabsResponse {
  result := SslLabsResponse{}

  var queryString = map[string]string{
    "host": domain,
  }

  body, err := apiClient.GetWithParams("/api/v3/analyze", queryString)

  if err != nil {
    log.Printf("Error making request with domain %s, err: %vx", domain, err)
  }

  json.Unmarshal(body, &result)
  result.AddExternalData(domain)

  return result
}

func (apiResponse *SslLabsResponse) AddExternalData(domain string) {
  for index, _ := range apiResponse.Servers {
    owner, country := WhoIs(apiResponse.Servers[index].Address)

    if apiResponse.Servers[index].Status == "Unable to connect to the server" {
      apiResponse.Servers[index].SslGrade = "U"      // Mark SslGrade as unknown if server is down
      apiResponse.IsDown = true
    }

    apiResponse.Servers[index].Country = country
    apiResponse.Servers[index].Owner = owner
  }

  title, logo := GetDomainHead(domain)
  apiResponse.Logo = logo
  apiResponse.Title = title
}

func WhoIs(ip string) (string, string) {
  command := fmt.Sprintf(WhoIsCmd, ip)
  out, err := exec.Command("bash", "-c", command).Output()

  if err != nil {
    log.Printf("Error executing WHOIS, err %v", err)
    return "", ""
  }

  trimmedOutput := strings.TrimRight(string(out), "\r\n")
  commandValues := strings.Split(trimmedOutput, " ")
  return commandValues[0], commandValues[1]
}
