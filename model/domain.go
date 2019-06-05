package model

import (
  "log"
  "server_id_api/db"
)

type Server struct {
  Address  string `json:"ipAddress"`
  SslGrade string `json:"grade"`
  Status   string `json:"statusMessage"`
  Country  string
  Owner    string
}

type Domain struct {
  Servers   []Server  `json:"endpoints"`
  Name      string
  Title     string
  Logo      string
  SslGrade  string
  IsDown    bool
}

const InsertDomainQuery = `
  INSERT INTO domains (name, ssl_grade, title, logo, is_down)
  VALUES ('%s', '%s', '%s', '%s', %t)
  RETURNING id;
`

const InsertServerQuery = `
  INSERT INTO servers (address, ssl_grade, status, country, owner, domain_id)
  VALUES ('%s', '%s','%s', '%s', '%s', %d)
  RETURNING id;
`

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
  }

  return id, err
}

func (domain *Domain) Persist() {
  id, err := domain.Insert()

  if err == nil {
    for _, server := range domain.Servers {
      server.Insert(id)
    }
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
