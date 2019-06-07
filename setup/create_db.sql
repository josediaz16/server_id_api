CREATE USER IF NOT EXISTS serverapi;
DROP DATABASE IF EXISTS $(DB_NAME);

CREATE DATABASE $(DB_NAME);

GRANT ALL ON DATABASE $(DB_NAME) TO serverapi;

CREATE TABLE $(DB_NAME).domains (
  id          SERIAL    PRIMARY KEY,
  name        STRING    NOT NULL,
  ssl_grade   STRING    NOT NULL,
  title       STRING    NOT NULL,
  logo        STRING    NOT NULL DEFAULT '',
  is_down     BOOL      NOT NULL,
  updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE $(DB_NAME).servers (
  id          SERIAL PRIMARY KEY,
  address     STRING NOT NULL,
  ssl_grade   STRING NOT NULL,
  status      STRING NOT NULL,
  country     STRING NOT NULL DEFAULT '',
  owner       STRING NOT NULL DEFAULT '',
  domain_id   INT    NOT NULL REFERENCES $(DB_NAME).domains (id) ON DELETE CASCADE,
  INDEX (domain_id)
);
