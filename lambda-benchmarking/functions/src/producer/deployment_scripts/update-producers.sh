#!/bin/bash
NAME=benchmarking
DIRECTORY_PATH=deployments
LOG_PATH="$DIRECTORY_PATH/updated_$(date '+%F_%H:%M:%S').txt"

GOOS=linux go build -v -race -o ./producer-handler handler.go
rm $NAME.zip
zip $NAME.zip producer-handler
mkdir -p $DIRECTORY_PATH

for ((deployIndex = $1; deployIndex < $2; deployIndex++)); do
  echo "Updating lambda producer $NAME-$deployIndex with newest code"
  /usr/local/bin/aws lambda update-function-code \
    --function-name "$NAME-$deployIndex" \
    --zip-file fileb://$NAME.zip >>"${LOG_PATH}"
done

echo "All producer Lambda functions from $1 to $2 were updated on AWS."
