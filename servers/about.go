package servers

import (
  "fmt"
  "log"
  "encoding/json"
  "os/exec"
  "strings"
  "server_id_api/api"
  "server_id_api/model"
  "sort"
)

const WhoIsCmd = "whois %s | grep -E \\(Country\\|OrgName\\) | awk '{print $2}' | xargs"

func GetServerData(apiClient *api.API, domain string) model.Domain {
  result := model.Domain{}

  var queryString = map[string]string{
    "host": domain,
  }

  body, err := apiClient.GetWithParams("/api/v3/analyze", queryString)

  if err != nil {
    log.Printf("Error making request with domain %s, err: %vx", domain, err)
  }

  json.Unmarshal(body, &result)
  addExternalData(domain, &result)
  result.Persist()

  return result
}

func addExternalData(domainName string, domain *model.Domain) {
  sslGrades := make([]string, len(domain.Servers))

  for index, _ := range domain.Servers {
    //owner, country := WhoIs(domain.Servers[index].Address)

    if domain.Servers[index].Status == "Unable to connect to the server" {
      domain.Servers[index].SslGrade = "U"      // Mark SslGrade as unknown if server is down
      domain.IsDown = true
    }

    sslGrades[index] = domain.Servers[index].SslGrade

    domain.Servers[index].Country = "US"
    domain.Servers[index].Owner = "Google"
  }

  title, logo := GetDomainHead(domainName)

  domain.Name = domainName
  domain.Logo = logo
  domain.Title = title
  domain.SslGrade = defineGlobalGrade(sslGrades)
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

func defineGlobalGrade(sslGrades []string) string {
  sort.Strings(sslGrades)
  return sslGrades[len(sslGrades) - 1]
}
