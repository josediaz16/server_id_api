package model

import (
  "log"
  "server_id_api/db"
  "time"
)

type Server struct {
  Address   string `json:"ipAddress"`
  SslGrade  string `json:"grade"`
  Status    string `json:"statusMessage"`
  Country   string
  Owner     string
  id        int
  domainId  int
}

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

const InsertServerQuery = `
  INSERT INTO servers (address, ssl_grade, status, country, owner, domain_id)
  VALUES ('%s', '%s','%s', '%s', '%s', %d)
  RETURNING id;
`

const TimeLayout = "2006-01-02 15:04:05 -0700"

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

func (domain *Domain) ShouldUpdate() bool {
  timestamp, _ := time.Parse(TimeLayout, domain.UpdatedAt)
  return timestamp.Before(time.Now().Add(-1*time.Hour))
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

  err := conn.FindBy("domains", "name", name).Scan(
    &domain.Id,
    &domain.Name,
    &domain.SslGrade,
    &domain.Title,
    &domain.Logo,
    &domain.IsDown,
    &domain.UpdatedAt,
  )

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

      err := rows.Scan(
        &server.id,
        &server.Address,
        &server.SslGrade,
        &server.Status,
        &server.Country,
        &server.Owner,
        &server.domainId,
      )

      if err == nil {
        servers = append(servers, server)
      }
    }
    return servers, nil
  }
}

func (server *Server) Insert(domainId int) (int, error) {
  conn := db.NewConn()

  id, err := conn.Insert(
    InsertServerQuery,
    server.Address,
    server.SslGrade,
    server.Status,
    server.Country,
    server.Owner,
    domainId,
  )

  if err != nil {
    log.Printf("Error Inserting server: %v", err)
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
