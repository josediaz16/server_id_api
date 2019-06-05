CREATE USER IF NOT EXISTS serverapi;
DROP DATABASE IF EXISTS servers_test;

CREATE DATABASE servers_test;

GRANT ALL ON DATABASE servers_test TO serverapi;

CREATE TABLE servers_test.domains (
  id          SERIAL    PRIMARY KEY,
  name        STRING    NOT NULL,
  ssl_grade   STRING    NOT NULL,
  title       STRING    NOT NULL,
  logo        STRING    NOT NULL DEFAULT '',
  updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE servers_test.servers (
  id          SERIAL PRIMARY KEY,
  address     STRING NOT NULL,
  ssl_grade   STRING NOT NULL,
  status      STRING NOT NULL,
  country     STRING NOT NULL DEFAULT '',
  owner       STRING NOT NULL DEFAULT '',
  domain_id   INT    NOT NULL REFERENCES servers_test.domains (id) ON DELETE CASCADE,
  INDEX (domain_id)
);
