import json
import boto3
from boto3.dynamodb.conditions import Key
import logging

import json
import logging
import os
import time
import uuid

dynamodb = boto3.resource('dynamodb')


def lambda_handler(event, context):
	logger = logging.getLogger()
	logger.setLevel(logging.INFO)
	logger.info("Request: %s", event)
	
	# Extracting query params
	experiment_type = event.get('queryStringParameters').get('experiment_type') # Mandatory
	start_date = event.get('queryStringParameters').get('start_date') # Mandatory
	end_date = event.get('queryStringParameters').get('end_date') # Optional
	
	if experiment_type is None:
		logging.error("Experiment type is required.")
		raise Exception("Couldn't fetch data. 'experiment_type' is required.")
        
	# DB configuration
	table = dynamodb.Table('continous_results')
	response = ''
	
	# Querying within a date range
	if(start_date and end_date):
		response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) & Key('date').between(start_date,end_date)
			)
	else:
		if(start_date):
			response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) & Key('date').gte(start_date)
			)
		else:
			response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) & Key('date').gte('2020-01-01')
			)
			
	items = response['Items']
	
	logger.info("Response: %s", items)
	# Create response
	return {
        "statusCode": 200,
        "body": json.dumps(items,default=str)
    }