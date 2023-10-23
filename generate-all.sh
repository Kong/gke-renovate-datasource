#!/bin/bash

set -euo pipefail

echo "Generating Stable channel"
go run . --channel stable --location "${GOOGLE_LOCATION}" > ./static/stable.json

echo "Generating Regular channel"
go run . --channel regular --location "${GOOGLE_LOCATION}" > ./static/regular.json

echo "Generating Rapid channel"
go run . --channel rapid --location "${GOOGLE_LOCATION}" > ./static/rapid.json
