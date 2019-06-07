#!/bin/bash

echo Wait for servers to be up
sleep 5

DATABASE="${DATABASE:-servers_dev}"

cat /create_db.sql | sed "s/\$(DB_NAME)/${DATABASE}/g" > /tmp/create_db.sql
cat /tmp/create_db.sql

/cockroach/cockroach.sh sql --host roach2 --insecure < /tmp/create_db.sql
