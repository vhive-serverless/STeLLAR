# MIT License
#
# Copyright (c) 2022 Dilina Dehigama and EASE Lab
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

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
	table = dynamodb.Table('continuous_results')
	response = ''
	
	# Querying within a date range
	if(start_date and end_date):
		response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) &
				 Key('date').between(start_date,end_date)
			)
	else:
		if(start_date):
			response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) &
				 Key('date').gte(start_date)
			)
		else:
			response = table.query(
				KeyConditionExpression=Key('experiment_type').eq(experiment_type) &
				 Key('date').gte('2020-01-01')
			)
			
	items = response['Items']
	
	logger.info("Response: %s", items)
	# Create response
	return {
        "statusCode": 200,
        "body": json.dumps(items,default=str)
    }