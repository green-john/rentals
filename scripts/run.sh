#!/usr/bin/env bash

export RENTALS_DB_HOST=localhost
export RENTALS_DB_NAME=rentals-test
export RENTALS_DB_USER=juan
./rentals-cli $1
