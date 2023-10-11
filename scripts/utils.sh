#! /bin/bash

load_env_vars() {
  if  [ -z "$ENV" ]; then
    source .env
  else
    source .env.$ENV
  fi
}
