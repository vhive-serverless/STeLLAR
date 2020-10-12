#!/bin/bash
NAME=benchmarking
REGION=eu-west-2
DIRECTORY_PATH=deployments
LOG_PATH="$DIRECTORY_PATH/created_$(date '+%F_%H:%M:%S').txt"
GATEWAYS_PATH="$DIRECTORY_PATH/gateways_$(date '+%F_%H:%M:%S').csv"
CLONE_API=c7532e9ynk
USAGE_PLAN_ID=af3hrd
STAGE=prod

GOOS=linux go build -v -race -o ./producer-handler handler.go
rm ${NAME}.zip
zip ${NAME}.zip producer-handler
mkdir -p ${DIRECTORY_PATH}

for ((deployIndex = $1; deployIndex < $2; deployIndex++)); do
  echo "Creating producer lambda: ${NAME}-${deployIndex}"
  /usr/local/bin/aws lambda create-function \
    --function-name "${NAME}-${deployIndex}" \
    --runtime go1.x \
    --role "$AWS_LAMBDA_ROLE" \
    --handler producer-handler \
    --zip-file fileb://$NAME.zip \
    --tracing-config Mode=PassThrough >>"${LOG_PATH}"
  # Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray. PassThrough otherwise.

  LAMBDAARN=$(/usr/local/bin/aws lambda list-functions \
    --query "Functions[?FunctionName==\`${NAME}-${deployIndex}\`].FunctionArn" \
    --output text \
    --region ${REGION})
  echo "ARN of lambda $NAME-$deployIndex is ${LAMBDAARN}"

  echo "Creating corresponding API: ${NAME}-API-${deployIndex} (clone of ${CLONE_API})"
  /usr/local/bin/aws apigateway create-rest-api \
    --name "${NAME}-API-$deployIndex" \
    --description "The API used to access benchmarking Lambda function $deployIndex." \
    --endpoint-configuration types=REGIONAL \
    --region ${REGION} \
    --clone-from ${CLONE_API} >>"${LOG_PATH}"

  APIID=$(/usr/local/bin/aws apigateway get-rest-apis \
    --query "items[?name==\`${NAME}-API-$deployIndex\`].id" \
    --output text)
  echo "API ID of ${NAME}-API-${deployIndex} is ${APIID}"

  RESOURCEID=$(/usr/local/bin/aws apigateway get-resources \
    --rest-api-id "${APIID}" \
    --query 'items[0].id' \
    --output text \
    --region ${REGION})
  echo "RESOURCEID of ${NAME}-API-${deployIndex} is ${RESOURCEID}"

  echo "Creating integration between lambda ${NAME}-${deployIndex} and API ${NAME}-API-${deployIndex}"
  /usr/local/bin/aws apigateway put-integration \
    --rest-api-id "${APIID}" \
    --resource-id "${RESOURCEID}" \
    --http-method ANY \
    --type AWS_PROXY \
    --integration-http-method ANY \
    --uri arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/${LAMBDAARN}/invocations \
    --request-templates '{"application/x-www-form-urlencoded":"{\"body\": $input.json(\"$\")}"}' \
    --region "${REGION}" >>"${LOG_PATH}"

  echo "Creating deployment for API ${NAME}-API-${deployIndex}, stage is ${STAGE}"
  /usr/local/bin/aws apigateway create-deployment \
    --rest-api-id "${APIID}" \
    --stage-name ${STAGE} \
    --region ${REGION} >>"${LOG_PATH}"

  echo "Adding API ID ${APIID} to general usage plan"
  /usr/local/bin/aws apigateway update-usage-plan \
    --usage-plan-id "${USAGE_PLAN_ID}" \
    --patch-operations "[{\"op\":\"add\",\"path\":\"/apiStages\",\"value\":\"${APIID}:${STAGE}\"}]" >>"${LOG_PATH}"

  echo "Adding permissions to execute lambda function ${NAME}-${deployIndex}"
  APIARN=$(echo "${LAMBDAARN}" | sed -e 's/lambda/execute-api/' -e "s/function:${NAME}-${deployIndex}/${APIID}/")
  /usr/local/bin/aws lambda add-permission \
    --function-name "${NAME}-${deployIndex}" \
    --statement-id apigateway-benchmarking \
    --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn "${APIARN}/prod/ANY/benchmarking" \
    --region ${REGION} >>"${LOG_PATH}"

  echo -n "${APIID} " >>"${GATEWAYS_PATH}"
done

echo "All producer Lambda functions from $1 to $2 were deployed to AWS."
