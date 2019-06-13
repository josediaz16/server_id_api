#!/bin/bash

DATABASE_URL="postgresql://serverapi@roach1:26257/servers_test?sslmode=disable"
go test ./test/... -v
