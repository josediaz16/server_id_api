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

const WhoIsCmd = "whois %s | grep -E \\(Country\\|OrgName\\) | awk -F ':' '{print $2}' | paste -sd \",\" | xargs"

const ServerReady = "READY"
const ServerDown = "Unable to connect to the server"
const UnkownStatus = "U"

func GetAllDomains() (map[string]*model.Domain, error) {
  var domainIds []int
  domainIdRegistry := make(map[int]string)
  domainRegistry := make(map[string]*model.Domain)
  domains, err := model.ListDomains()

  if err != nil {
    return domainRegistry, err
  }

  for index, _ := range domains {
    domain := domains[index]
    domainIds = append(domainIds, domain.Id)
    domainIdRegistry[domain.Id] = domain.Name
    domainRegistry[domain.Name] = &domain
  }

  if len(domainIds) == 0 {
    return domainRegistry, err
  }

  servers, _ := model.GetRelatedServers(domainIds)

  for servers.Next() {
    var server model.Server
    server.FromDb(servers)
    domainName := domainIdRegistry[server.GetDomainId()]
    domainRegistry[domainName].Servers = append(domainRegistry[domainName].Servers, server)
  }

  return domainRegistry, nil
}

func GetServerData(apiClient *api.API, domain string) model.Domain {
  var result model.Domain

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
    owner, country := WhoIs(domain.Servers[index].Address)

    if domain.Servers[index].Status == ServerDown {
      domain.Servers[index].SslGrade = UnkownStatus      // Mark SslGrade as unknown if server is down
      domain.IsDown = true
    }

    sslGrades[index] = domain.Servers[index].SslGrade

    domain.Servers[index].Country = country
    domain.Servers[index].Owner = owner
  }

  title, logo := GetDomainHead(domainName)

  domain.Name = domainName
  domain.Logo = logo
  domain.Title = title
  domain.IsDown = domain.IsDown || domain.Status != ServerReady
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
  commandValues := strings.Split(trimmedOutput, ",")
  return commandValues[0], strings.TrimSpace(commandValues[1])
}

func defineGlobalGrade(sslGrades []string) string {
  sort.Strings(sslGrades)
  numberOfGrades := len(sslGrades)

  if numberOfGrades != 0 {
    return sslGrades[len(sslGrades) - 1]
  }

  return UnkownStatus
}
