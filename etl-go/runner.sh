#!/bin/bash

# Build etl-runner program if doesn't exist
if [[ ! -f "./etl-runner" ]]; then
  echo "Building etl-runner program"
  go build etl-runner.go
else
  echo "etl-runner already built. Skipping build"
fi

# Run etl-runner program
./etl-runner