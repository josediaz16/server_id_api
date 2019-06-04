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
  Endpoints []Server  `json:"endpoints"`
  Host      string    `json:"host"`
  Port      int       `json:"port"`
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
  buildServerData(&result)

  return result
}

func buildServerData(apiResponse *SslLabsResponse) {
  for index, _ := range apiResponse.Endpoints {
    owner, country := WhoIs(apiResponse.Endpoints[index].Address)

    apiResponse.Endpoints[index].Country = country
    apiResponse.Endpoints[index].Owner = owner
  }
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
