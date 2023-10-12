#!/bin/bash

echo "Generating Stable channel"
go run . --channel stable  > ./static/stable.json

echo "Generating Regular channel"
go run . --channel regular > ./static/regular.json

echo "Generating Rapid channel"
go run . --channel rapid > ./static/rapid.json
