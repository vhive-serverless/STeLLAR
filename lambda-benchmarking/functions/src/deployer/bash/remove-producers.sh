#!/bin/bash
NAME=benchmarking

for ((deployIndex = $1; deployIndex < $2; deployIndex++)); do
  echo "Removing producer ${NAME}-$deployIndex"
  /usr/local/bin/aws lambda delete-function \
    --function-name "${NAME}-$deployIndex"

  APIID=$(/usr/local/bin/aws apigateway get-rest-apis \
    --query "items[?name==\`${NAME}-API-$deployIndex\`].id" \
    --output text)
  echo "API ID of ${NAME}-API-${deployIndex} is ${APIID}"

  echo "Removing API ${NAME}-API-$deployIndex with API ID ${APIID}"
  /usr/local/bin/aws apigateway delete-rest-api \
    --rest-api-id "${APIID}"
done

echo "All producer Lambda functions from $1 to $2 were removed from AWS!"