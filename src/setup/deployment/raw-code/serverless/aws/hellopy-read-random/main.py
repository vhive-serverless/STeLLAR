import json
import os
import time
import random


def lambda_handler(request, context):
    incr_limit = 0

    if 'queryStringParameters' in request and 'IncrementLimit' in request['queryStringParameters']:
        incr_limit = int(request['queryStringParameters'].get('IncrementLimit', 0))
    elif 'body' in request and json.loads(request['body'])['IncrementLimit']:
        incr_limit = int(json.loads(request['body'])['IncrementLimit'])

    simulate_work(incr_limit)
    read_filler_file('./filler.file')

    json_region = os.environ.get('AWS_REGION', 'Unknown')

    response = {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "Region ": json_region,
            "RequestID": context.aws_request_id,
            "TimestampChain": [str(time.time_ns())]
        }, indent=4)
    }

    return response


def simulate_work(increment):
    # MAXNUM = 6103705
    num = 0
    while num < increment:
        num += 1


def read_filler_file(path: str) -> None:
    file_size = os.stat(path).st_size
    number_of_pages = file_size // 4096
    with open(path, 'rb') as f:
        for _ in range(100):
            page_number = random.randrange(0, number_of_pages)
            page_offset = page_number * 4096
            f.seek(page_offset)
            f.read(1)
