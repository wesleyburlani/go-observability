#!/bin/bash

local_dir=$(dirname "$(readlink -f "$0")")
root_dir="$local_dir/.."
source "$local_dir/utils.sh"

load_env_vars

echo "Running migrations up on $ENV environment..."
docker run -v "$root_dir/sql/migrations":/migrations --network host migrate/migrate -path=/migrations/ -database $DATABASE_URL up
