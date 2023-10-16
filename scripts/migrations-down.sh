#!/bin/bash

local_dir=$(dirname "$(readlink -f "$0")")
root_dir="$local_dir/.."
source "$local_dir/utils.sh"

load_env_vars

echo "Running migrations down on $ENV environment..."
docker run -v "$root_dir/sql/migrations":/migrations --network go-api_app_net migrate/migrate -path=/migrations/ -database $DATABASE_URL down -all
