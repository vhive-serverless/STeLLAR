import json

import time


def hello_world(request):
    request_json = request.get_json()

    incr_limit = 0
    if request.args and 'incrementLimit' in request.args:
        incr_limit = request.args.get('incrementLimit')
    elif request_json and 'incrementLimit' in request_json:
        incr_limit = request_json['incrementLimit']

    simulate_work(incr_limit)

    response = {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": {
            # "Region ": json_region,
            "RequestID": "google-does-not-specify",
            "TimestampChain": [str(time.time_ns())],
        }
    }

    return json.dumps(response, indent=4)


def simulate_work(incr):
    # MAXNUM = 6103705
    num = 0
    while num < incr:
        num += 1
