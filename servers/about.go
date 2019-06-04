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

func GetServerData(domain string) SslLabsResponse {
  client := &http.Client{}

  req, _ := http.NewRequest("GET", "https://api.ssllabs.com/api/v3/analyze", nil)
  query := req.URL.Query()

  query.Add("host", domain)
  req.URL.RawQuery = query.Encode()

  resp, err := client.Do(req)

  if err != nil {
    log.Printf("Error making request with domain %s, err: %vx", domain, err)
  }

  defer resp.Body.Close()

  result := SslLabsResponse{}

  temp, _ := ioutil.ReadAll(resp.Body)

  err = json.Unmarshal(temp, &result)
  buildServerData(&result)

  return result
}

func buildServerData(apiResponse *SslLabsResponse) {
  for index, _ := range apiResponse.Endpoints {
    command := fmt.Sprintf("whois %s | grep -E \\(Country\\|OrgName\\) | awk '{print $2}' | xargs", apiResponse.Endpoints[index].Address)

    out, _ := exec.Command("bash", "-c", command).Output()
    commandValues := strings.Split(strings.TrimRight(string(out), "\r\n"), " ")

    log.Printf("command %v", commandValues)
    apiResponse.Endpoints[index].Country = commandValues[1]
    apiResponse.Endpoints[index].Owner = commandValues[0]
  }
}
