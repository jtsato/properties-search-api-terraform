#!/bin/bash

# Build etl-go program if doesn't exist
if [[ ! -f "./etl-go" ]]; then
  echo "Building etl-go program"
  go build
else
  echo "Skipping build because etl-go program already exists"
fi

# Set permissions to execute etl-go program
sudo chmod +x etl-go

# Run etl-go program
./etl-go
