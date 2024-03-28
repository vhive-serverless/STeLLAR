import json
import time


def main(event, context):
    event = json.loads(event)

    increment_limit = 0
    if "IncrementLimit" in event["queryParameters"]:
        increment_limit = int(event["queryParameters"]["IncrementLimit"])
    if "IncrementLimit" in event["body"]:
        increment_limit = int(event["body"]["IncrementLimit"])
    simulate_work(increment_limit)

    response_body = {
        "Region": context.region,
        "RequestID": context.request_id,
        "TimestampChain": [str(time.time_ns())],
    }
    response = {
        "isBase64Encoded": "false",
        "statusCode": "200",
        "headers": {"x-custom-header": "no", "Content-Type": "application/json"},
        "body": response_body,
    }
    return json.dumps(response)


def simulate_work(increment_limit: int) -> None:
    num = 0
    while num < increment_limit:
        num += 1
