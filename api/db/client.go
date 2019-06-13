package db

import (
  "database/sql"
  "log"
  "fmt"
  "os"
  "encoding/json"
  "strings"
  _ "github.com/lib/pq"
)

type Conn struct {
  Client *sql.DB
}

func NewConn() (*Conn) {
  db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

  if err != nil {
    log.Fatal("error connecting to database: ", err)
  }

  return &Conn{db}
}

// PERSIST OPERATIONS

func (conn *Conn) Insert(query string, args ...interface{}) (int, error) {
  return conn.Persist("Success Insert", query, args...)
}

func (conn *Conn) Update(query string, args ...interface{}) (int, error) {
  return conn.Persist("Success Update", query, args...)
}

func (conn *Conn) Persist(logMessage string, query string, args ...interface{}) (int, error) {
  var id int
  formattedQuery := fmt.Sprintf(query, args...)

  err := conn.Client.QueryRow(formattedQuery).Scan(&id)

  if err != nil {
    return 0, err
  } else {
    log.Printf("%s: %d", logMessage, id)
    return id, nil
  }
}

// SELECT OPERATIONS

func (conn *Conn) Count(tableName string) (count int, err error) {
  formattedQuery := fmt.Sprintf("SELECT COUNT(*) as count FROM %s;", tableName)

  err = conn.Client.QueryRow(formattedQuery).Scan(&count)
  return count, err
}

func (conn *Conn) FindBy(tableName string, field string, value interface{}) (*sql.Row) {
  formattedQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1;", tableName, field)
  return conn.Client.QueryRow(formattedQuery, value)
}

func (conn * Conn) GetAll(tableName string) (*sql.Rows, error) {
  formattedQuery := fmt.Sprintf("SELECT * FROM %s", tableName)
  return conn.Client.Query(formattedQuery)
}

func (conn *Conn) GetChildRecords(tableName string, foreignKey string, id int) (*sql.Rows, error) {
  formattedQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1;", tableName, foreignKey)
  return conn.Client.Query(formattedQuery, id)
}

func (conn *Conn) GetAllChilds(tableName string, foreignKey string, ids []int) (*sql.Rows, error) {
  bytes, _ := json.Marshal(ids)
  strIds := strings.Trim(string(bytes), "[]")

  formattedQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s IN(%s);", tableName, foreignKey, strIds)
  return conn.Client.Query(formattedQuery)
}

// DELETE OPERATIONS

func (conn *Conn) DeleteChildRecords(tableName string, foreignKey string, id int) (count int, err error) {
  formattedQuery := fmt.Sprintf("DELETE FROM %s WHERE %s = %d;", tableName, foreignKey, id)
  return conn.Delete(formattedQuery)
}

func (conn *Conn) DeleteAllFrom(tableName string) (int, error) {
  formattedQuery := fmt.Sprintf("DELETE FROM %s;", tableName)
  return conn.Delete(formattedQuery)
}

func (conn *Conn) Delete(query string) (count int, err error) {
  result, err := conn.Client.Exec(query)

  if err != nil {
    return 0, err
  } else {
    rows, _ := result.RowsAffected()
    log.Printf("Success Delete Rows: %d", rows)
    return int(rows), err
  }
}
