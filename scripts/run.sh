#!/bin/sh

set -euo pipefail

CURRENT_DIR="$(cd $(dirname "$0") && pwd)"

sh "$CURRENT_DIR/build.sh"

docker build -t boringdownload/drone-junit -f "$CURRENT_DIR/../docker/Dockerfile" "$CURRENT_DIR/../"

docker run --rm -e PLUGIN_PATHS="report.xml" -e PLUGIN_REPORT_NAME="drone-junit" \
  -e PLUGIN_LOG_LEVEL="debug" \
  -e DRONE_CARD_PATH="/dev/stdout" \
  -e DRONE_COMMIT_SHA=8f51ad7884c5eb69c11d260a31da7a745e6b78e2 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_BUILD_NUMBER=43 \
  -e DRONE_BUILD_STATUS=success \
  -w /drone/src \
  -v "$CURRENT_DIR/../:/drone/src" \
  boringdownload/drone-junit
