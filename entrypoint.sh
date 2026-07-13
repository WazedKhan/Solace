#!/bin/sh
set -e
migrate -path /app/migrations/ -database "$DSN" up
exec ./main
