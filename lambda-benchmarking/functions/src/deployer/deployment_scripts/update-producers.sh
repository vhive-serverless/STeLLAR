#!/bin/bash
NAME=benchmarking
DIRECTORY_PATH=deployments
LOG_PATH="$DIRECTORY_PATH/updated_$(date '+%F_%H:%M:%S').txt"

echo "Updating lambda producer $NAME-$1 with newest code"
/usr/local/bin/aws lambda update-function-code \
  --function-name "$NAME-$1" \
  --zip-file fileb://$NAME.zip >>"${LOG_PATH}"