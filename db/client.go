package db

import (
  "database/sql"
  "log"
  "fmt"
  _ "github.com/lib/pq"
)

type Conn struct {
  Client *sql.DB
}

func NewConn() (*Conn) {
  db, err := sql.Open("postgres", "postgresql://serverapi@roach1:26257/servers_test?sslmode=disable")

  if err != nil {
    log.Fatal("error connecting to database: ", err)
  }

  return &Conn{db}
}

func (conn *Conn) Insert(query string, args ...interface{}) (int, error) {
  var id int
  formattedQuery := fmt.Sprintf(query, args...)

  err := conn.Client.QueryRow(formattedQuery).Scan(&id)

  if err != nil {
    return 0, err
  } else {
    log.Printf("Success Insert: %d", id)
    return id, nil
  }
}

func (conn *Conn) DeleteAllFrom(tableName string) (int, error) {
  formattedQuery := fmt.Sprintf("DELETE FROM %s;", tableName)

  result, err := conn.Client.Exec(formattedQuery)

  if err != nil {
    return 0, err
  } else {
    rows, _ := result.RowsAffected()
    log.Printf("Success Delete Rows: %d", rows)
    return int(rows), err
  }
}

func (conn *Conn) Count(tableName string) (count int, err error) {
  formattedQuery := fmt.Sprintf("SELECT COUNT(*) as count FROM %s;", tableName)

  err = conn.Client.QueryRow(formattedQuery).Scan(&count)
  return count, err
}
