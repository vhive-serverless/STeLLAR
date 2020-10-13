#!/bin/bash
NAME=benchmarking

echo "Removing producer ${NAME}-$1"
/usr/local/bin/aws lambda delete-function --function-name "${NAME}-$1"

APIID=$(/usr/local/bin/aws apigateway get-rest-apis \
  --query "items[?name==\`${NAME}-API-$1\`].id" \
  --output text)
echo "API ID of ${NAME}-API-$1 is ${APIID}"

echo "Removing API ${NAME}-API-$1 with API ID ${APIID}"
/usr/local/bin/aws apigateway delete-rest-api --rest-api-id "${APIID}"
