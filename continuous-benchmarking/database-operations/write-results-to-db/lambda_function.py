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

    table = dynamodb.Table('continuous_results')

    # Create a table entry
    entry = {
        'id': str(uuid.uuid1()),
        'experiment_type': experiment_type,
        'subtype': data.get('subtype'),
        'min': data.get('min'),
        'max': data.get('max'),
        'median': data.get('median'),
        'tail_latency': data.get('tail_latency'),
        'first_quartile': data.get('first_quartile'),
        'third_quartile': data.get('third_quartile'),
        'standard_deviation': data.get('standard_deviation'),
        'payload_size': data.get('payload_size'),
        'burst_size': data.get('burst_size'),
        'IAT_type': data.get('IAT_type'),
        'count': data.get('count'),
        'date': date,
        'provider': data.get('provider'),
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
    
