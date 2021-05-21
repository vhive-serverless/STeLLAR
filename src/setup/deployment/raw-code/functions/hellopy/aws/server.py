from __future__ import print_function

import json
import os


def handler(_, context):
    json_region = os.environ['AWS_REGION']
    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "Region ": json_region,
            "RequestID": context.aws_request_id,
            "TimestampChain": '[0]',
        })
    }
