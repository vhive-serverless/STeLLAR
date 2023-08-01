from __future__ import print_function

import json
import os
import time


def lambda_handler(request, context):

    incr_limit = 0

    if(request['queryStringParameters'] and 'IncrementLimit' in request['queryStringParameters']):
        incr_limit = int(request['queryStringParameters'].get('IncrementLimit', 0))
    elif request['body'] and json.loads(request['body'])['IncrementLimit']:
        incr_limit = int(json.loads(request['body'])['IncrementLimit'])

    simulate_work(incr_limit)

    json_region = os.environ.get('AWS_REGION','Unknown')

    response = {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "Region ": json_region,
            "RequestID": context.aws_request_id,
            "TimestampChain": [str(time.time_ns())]
        },indent=4)
    }

    return response


def simulate_work(increment):
    # MAXNUM = 6103705
    num = 0
    while num < increment:
        num += 1
