#!/usr/bin/env bash

cd ..
export RENTALS_DB_HOST=localhost
export RENTALS_DB_NAME=rentals-dev
export RENTALS_DB_USER=juan
go test ./...

