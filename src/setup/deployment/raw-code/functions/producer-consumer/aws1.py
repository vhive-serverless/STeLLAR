import os
import boto3
import boto3.session
import json

sessionInstance = None

def invokeNextFunctionAWS(paramters, functionID):
	namingPrefix = 'vHive-bench_'
	
	nextFunctionPayload = json.dumps(parameters)
	print(f"Invoking next function: {namingPrefix}{functionID}")

	lambdaClient = autheticateLambdaClient()

	response = lambdaClient.invoke(
        FunctionName=f"{namingPrefix}{functionID}",
        InvocationType="RequestResponse",
		LogType="Tail",
        Payload=bytes(nextFunctionPayload, encoding="utf-8"),
    )

	# Steaming body, .read() or json.load() if expecting json response
	return response["Payload"]

def createSessionInstance():

	region = os.environ['AWS_REGION']
	createdSessionInstance = boto3.session.Session(region_name=region)
	return createdSessionInstance

def autheticateLambdaClient():

	if sessionInstance is None:
		sessionInstance = createSessionInstance()


	return sessionInstance.client('lambda')

def autheticates3Client():

	if sessionInstance is None:
		sessionInstance = createSessionInstance()


	return sessionInstance.client('s3') #add s3 uploader
	
