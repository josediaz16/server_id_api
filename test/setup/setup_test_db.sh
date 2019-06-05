#!/bin/bash

echo Wait for servers to be up
sleep 5

/cockroach/cockroach.sh sql --host roach2 --insecure < /create_db.sql
