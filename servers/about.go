package servers

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "encoding/json"
  "os/exec"
  "strings"
)

type Server struct {
  Address  string `json:"ipAddress"`
  SslGrade string `json:"grade"`
  Country  string
  Owner    string
}

type SslLabsResponse struct {
  Servers   []Server  `json:"endpoints"`
  Title     string
  Logo      string
}

const WhoIsCmd = "whois %s | grep -E \\(Country\\|OrgName\\) | awk '{print $2}' | xargs"

func GetServerData(domain string) SslLabsResponse {
  client := &http.Client{}
  result := SslLabsResponse{}

  req, _ := http.NewRequest("GET", "https://api.ssllabs.com/api/v3/analyze", nil)
  query := req.URL.Query()

  query.Add("host", domain)
  req.URL.RawQuery = query.Encode()

  resp, err := client.Do(req)

  if err != nil {
    log.Printf("Error making request with domain %s, err: %vx", domain, err)
  }

  defer resp.Body.Close()

  temp, _ := ioutil.ReadAll(resp.Body)

  err = json.Unmarshal(temp, &result)
  result.AddExternalData(domain)

  return result
}

func (apiResponse *SslLabsResponse) AddExternalData(domain string) {
  for index, _ := range apiResponse.Servers {
    owner, country := WhoIs(apiResponse.Servers[index].Address)

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
