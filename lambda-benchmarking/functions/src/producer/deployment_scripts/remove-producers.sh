#!/bin/bash
NAME=benchmarking

for ((deployIndex = 0; deployIndex < $1; deployIndex++)); do
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

echo "All $1 producer Lambda functions were removed from AWS!"