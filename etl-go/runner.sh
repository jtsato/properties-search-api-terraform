#!/bin/bash

# Build etl-go program if doesn't exist
if [[ ! -f "./etl-go" ]]; then
  echo "Building etl-go program"
  go build
  sudo chmod +x etl-go
else
  echo "Skipping build because etl-go program already exists"
fi

# Run etl-go program
./etl-go
