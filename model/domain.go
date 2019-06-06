package model

import (
  "log"
  "server_id_api/db"
  "time"
  "os"
  "strconv"
)

type Domain struct {
  Servers           []Server  `json:"endpoints"`
  Status            string    `json:"status"`
  Id                int
  Name              string
  Title             string
  Logo              string
  SslGrade          string
  PreviousSslGrade  string
  IsDown            bool
  ServersChanged    bool
  UpdatedAt         string
}

type Row interface {
  Scan(dest ...interface{}) error
}

const InsertDomainQuery = `
  INSERT INTO domains (name, ssl_grade, title, logo, is_down)
  VALUES ('%s', '%s', '%s', '%s', %t)
  RETURNING id;
`

const UpdateDomainQuery = `
  UPDATE domains set(ssl_grade, title, logo, is_down, updated_at) = ('%s', '%s', '%s', %t, '%s')
  WHERE id = %d
  RETURNING id;
`

const TimeLayout = "2006-01-02 15:04:05 -0700"

func TimeWindow() time.Time {
  prevTime := os.Getenv("UPDATE_TIME_WINDOW")

  if prevTime == "" {
    prevTime = "1"
  }

  hours, _ := strconv.Atoi(prevTime)
  return time.Now().Add(-time.Duration(hours) * time.Hour)
}

func (domain *Domain) Insert() (int, error) {
  conn := db.NewConn()

  id, err := conn.Insert(
    InsertDomainQuery,
    domain.Name,
    domain.SslGrade,
    domain.Title,
    domain.Logo,
    domain.IsDown,
  )

  if err != nil {
    log.Printf("Error Inserting domain: %v", err)
  } else {
    domain.Id = id
  }

  return id, err
}

func (domain *Domain) Update() (int, error) {
  conn := db.NewConn()

  id, err := conn.Update(
    UpdateDomainQuery,
    domain.SslGrade,
    domain.Title,
    domain.Logo,
    domain.IsDown,
    domain.UpdatedAt,
    domain.Id,
  )

  if err != nil {
    log.Printf("Error Updating Domain: %v", err)
  }

  return id, err
}

func (domain *Domain) FromDb(dataset Row) error {
  err := dataset.Scan(
    &domain.Id,
    &domain.Name,
    &domain.SslGrade,
    &domain.Title,
    &domain.Logo,
    &domain.IsDown,
    &domain.UpdatedAt,
  )
  return err
}

func (domain *Domain) Persist() {
  prevDomain, _ := GetDomainByName(domain.Name)

  if prevDomain != nil {
    if prevDomain.ShouldUpdate() {
      domain.Id = prevDomain.Id
      domain.UpdatedAt = time.Now().Format(TimeLayout)
      domain.UpdateWithServers()
    }

    domain.ServersChanged = domain.SslGrade != prevDomain.SslGrade
    domain.PreviousSslGrade = prevDomain.SslGrade
  } else {
    domain.InsertWithServers()
  }
}

func ListDomains() ([]Domain, error) {
  var domains []Domain
  var domainIds []int

  conn := db.NewConn()
  rows, err := conn.GetAll("domains")

  defer rows.Close()

  if err != nil {
    return domains, err
  } else {

    domainRegistry := make(map[int]*Domain)

    for rows.Next() {
      var domain Domain
      domain.FromDb(rows)

      domainIds = append(domainIds, domain.Id)
      domainRegistry[domain.Id] = &domain
    }

    if len(domainIds) > 0 {
      servers, _ := conn.GetAllChilds("servers", "domain_id", domainIds)

      for servers.Next() {
        var server Server
        server.FromDb(servers)
        domainRegistry[server.domainId].Servers = append(domainRegistry[server.domainId].Servers, server)
      }
    }

    for _, domain := range domainRegistry {
      domains = append(domains, *domain)
    }

    return domains, nil
  }
}

func (domain *Domain) ShouldUpdate() bool {
  timestamp, _ := time.Parse(TimeLayout, domain.UpdatedAt)
  return timestamp.Before(TimeWindow())
}

func (domain *Domain) InsertWithServers() {
  id, err := domain.Insert()

  if err == nil {
    for _, server := range domain.Servers {
      server.Insert(id)
    }
  }
}

func (domain *Domain) UpdateWithServers() {
  domain.DeleteServers()
  id, err := domain.Update()

  if err == nil {
    for _, server := range domain.Servers {
      server.Insert(id)
    }
  }
}

func (domain *Domain) DeleteServers() {
  conn := db.NewConn()
  conn.DeleteChildRecords("servers", "domain_id", domain.Id)
}

func GetDomainByName(name string) (*Domain, error) {
  conn := db.NewConn()
  var domain Domain

  err := domain.FromDb(conn.FindBy("domains", "name", name))

  if err != nil {
    return nil, err
  } else {
    servers, _ := GetServersByDomain(domain.Id)
    domain.Servers = servers
    return &domain, nil
  }
}

func GetServersByDomain(domainId int) ([]Server, error) {
  var servers []Server

  conn := db.NewConn()
  rows, err := conn.GetChildRecords("servers", "domain_id", domainId)

  defer rows.Close()

  if err != nil {
    return servers, err
  } else {

    for rows.Next() {
      var server Server
      server.FromDb(rows)
      servers = append(servers, server)
    }

    return servers, nil
  }
}
