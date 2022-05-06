#! /bin/bash

tries=0
until [ "$tries" -ge 5 ]; do
  docker-compose exec -T mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "p@ssw0rd" -Q "SELECT 1" && break
  tries=$((tries+1))
  sleep 2
done
