#! /bin/bash

tries=0
until [ "$tries" -ge 5 ]; do
  docker-compose exec -T mssql /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "p@ssw0rd" -No -Q "SELECT 1" && break
  tries=$((tries+1))
  sleep 2
done
