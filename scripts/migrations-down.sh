#!/bin/bash

local_dir=$(dirname "$(readlink -f "$0")")
root_dir="$local_dir/.."
source "$local_dir/utils.sh"

load_env_vars

echo "Running migrations down on $ENV environment..."
migrate -source "file://./sql/migrations" -database $DATABASE_URL down -all
