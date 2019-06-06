package model

import (
  "log"
  "server_id_api/db"
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

const InsertServerQuery = `
  INSERT INTO servers (address, ssl_grade, status, country, owner, domain_id)
  VALUES ('%s', '%s','%s', '%s', '%s', %d)
  RETURNING id;
`

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

func (server *Server) FromDb(dataset Row) error {
  err := dataset.Scan(
    &server.id,
    &server.Address,
    &server.SslGrade,
    &server.Status,
    &server.Country,
    &server.Owner,
    &server.domainId,
  )

  return err
}

