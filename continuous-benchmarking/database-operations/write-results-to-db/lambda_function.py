import json
import logging
import os
import time
import uuid

import boto3
dynamodb = boto3.resource('dynamodb')

def lambda_handler(event, context):
    
    logger = logging.getLogger()
    logger.setLevel(logging.INFO)
    logger.info("Request: %s", event)
	
    data = json.loads(str(event['body']).replace("'",'"'))
    response = ''
    
    # Partition key - experiment_type (Mandatory)
    # Sort key - date (Mandatory)
    
    experiment_type = data.get('experiment_type')
    date = data.get('date')
    
    if not (experiment_type and date):
        response = 'required exp or date'
        logging.error("Experiment type is required.")
        raise Exception("Couldn't write data.")
        
    timestamp = str(time.time())

    table = dynamodb.Table('continous_results')
    
    # Extract params from body
    subtype = data.get('subtype')
    minimum = data.get('min')
    maximum = data.get('max')
    median = data.get('median')
    tail_latency = data.get('tail_latency')
    first_quartile = data.get('first_quartile')
    third_quartile = data.get('third_quartile')
    standard_deviation = data.get('standard_deviation')
    payload_size = data.get('payload_size')
    burst_size = data.get('burst_size')
    IAT_type = data.get('IAT_type')
    count = data.get('count')
    date = data.get('date')
    provider = data.get('provider')

    # Create a table entry
    entry = {
        'id': str(uuid.uuid1()),
        'experiment_type': experiment_type,
        'subtype': subtype,
        'min': minimum,
        'max': maximum,
        'median': median,
        'tail_latency': tail_latency,
        'first_quartile': first_quartile,
        'third_quartile': third_quartile,
        'standard_deviation': standard_deviation,
        'payload_size': payload_size,
        'burst_size': burst_size,
        'IAT_type': IAT_type,
        'count': count,
        'date': date,
        'provider': provider,
        'createdAt': timestamp
    }

    # write the entry to the database
    table.put_item(Item=entry)

    # create a response
    response = {
        "statusCode": 200,
        "body": json.dumps(entry)
    }

    return response
    
