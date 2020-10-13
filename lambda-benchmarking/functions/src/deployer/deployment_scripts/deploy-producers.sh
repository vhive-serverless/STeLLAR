#!/bin/bash
NAME=benchmarking
REGION=us-east-2
DIRECTORY_PATH=deployments
LOG_PATH="$DIRECTORY_PATH/created_$(date '+%F_%H:%M:%S').txt"
GATEWAYS_PATH="$DIRECTORY_PATH/gateways_$(date '+%F_%H:%M:%S').csv"
CLONE_API=i6lttb52fj
USAGE_PLAN_ID=gdkx9z
STAGE=prod

echo "Creating producer lambda: ${NAME}-$1"
/usr/local/bin/aws lambda create-function \
  --function-name "${NAME}-$1" \
  --runtime go1.x \
  --role "$AWS_LAMBDA_ROLE" \
  --handler producer-handler \
  --zip-file fileb://../producer/$NAME.zip \
  --tracing-config Mode=PassThrough >>"${LOG_PATH}"
# Set Mode to Active to sample and trace a subset of incoming requests with AWS X-Ray. PassThrough otherwise.

LAMBDAARN=$(/usr/local/bin/aws lambda list-functions \
  --query "Functions[?FunctionName==\`${NAME}-$1\`].FunctionArn" \
  --output text \
  --region ${REGION})
echo "ARN of lambda $NAME-$1 is ${LAMBDAARN}"

echo "Creating corresponding API: ${NAME}-API-$1 (clone of ${CLONE_API})"
/usr/local/bin/aws apigateway create-rest-api \
  --name "${NAME}-API-$1" \
  --description "The API used to access benchmarking Lambda function $1." \
  --endpoint-configuration types=REGIONAL \
  --region ${REGION} \
  --clone-from ${CLONE_API} >>"${LOG_PATH}"

APIID=$(/usr/local/bin/aws apigateway get-rest-apis \
  --query "items[?name==\`${NAME}-API-$1\`].id" \
  --output text \
  --region ${REGION})
echo "API ID of ${NAME}-API-$1 is ${APIID}"

RESOURCEID=$(/usr/local/bin/aws apigateway get-resources \
  --rest-api-id "${APIID}" \
  --query 'items[1].id' \
  --output text \
  --region ${REGION})
echo "RESOURCEID of ${NAME}-API-$1 is ${RESOURCEID}"

echo "Creating integration between lambda ${NAME}-$1 and API ${NAME}-API-$1"
/usr/local/bin/aws apigateway put-integration \
  --rest-api-id "${APIID}" \
  --resource-id "${RESOURCEID}" \
  --http-method ANY \
  --type AWS_PROXY \
  --integration-http-method ANY \
  --uri arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/"${LAMBDAARN}"/invocations \
  --request-templates '{"application/x-www-form-urlencoded":"{\"body\": $input.json(\"$\")}"}' \
  --region "${REGION}" >>"${LOG_PATH}"

echo "Creating deployment for API ${NAME}-API-$1, stage is ${STAGE}"
/usr/local/bin/aws apigateway create-deployment \
  --rest-api-id "${APIID}" \
  --stage-name ${STAGE} \
  --region ${REGION} >>"${LOG_PATH}"

echo "Adding API ID ${APIID} to general usage plan"
/usr/local/bin/aws apigateway update-usage-plan \
  --usage-plan-id "${USAGE_PLAN_ID}" \
  --patch-operations "[{\"op\":\"add\",\"path\":\"/apiStages\",\"value\":\"${APIID}:${STAGE}\"}]" >>"${LOG_PATH}"

echo "Adding permissions to execute lambda function ${NAME}-$1"
APIARN=$(echo "${LAMBDAARN}" | sed -e 's/lambda/execute-api/' -e "s/function:${NAME}-$1/${APIID}/")
/usr/local/bin/aws lambda add-permission \
  --function-name "${NAME}-$1" \
  --statement-id apigateway-benchmarking \
  --action lambda:InvokeFunction \
  --principal apigateway.amazonaws.com \
  --source-arn "${APIARN}/prod/ANY/benchmarking" \
  --region ${REGION} >>"${LOG_PATH}"

echo "${APIID} " >>"${GATEWAYS_PATH}"
