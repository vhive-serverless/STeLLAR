from __future__ import print_function

import json
import os
import time


def lambda_handler(request, context):
    request_json = request.get_json()

    incr_limit = 0
    if request.args and 'IncrementLimit' in request.args:
        incr_limit = request.args.get('IncrementLimit')
    elif request_json and 'IncrementLimit' in request_json:
        incr_limit = request_json['IncrementLimit']

    simulate_work(incr_limit)

    json_region = os.environ['AWS_REGION']
    response = {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": {
            "Region ": json_region,
            "RequestID": context.aws_request_id,
            "TimestampChain": [str(time.time_ns())],
        }
    }

    return json.dumps(response, indent=4)


def simulate_work(increment):
    # MAXNUM = 6103705
    num = 0
    while num < increment:
        num += 1
