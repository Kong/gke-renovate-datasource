#!/bin/bash

set -euo pipefail

generate_channel() {
	local channel=$1
	echo "Generating ${channel} channel"
	go run . --channel "${channel}" --location "${GOOGLE_LOCATION}" --out "./static/${channel}.json"
}

generate_channel "stable"
generate_channel "regular"
generate_channel "rapid"
