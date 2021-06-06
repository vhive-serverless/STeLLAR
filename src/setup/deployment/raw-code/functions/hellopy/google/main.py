import json
import time


def hello_world(request):
    """Responds to any HTTP request.
    Args:
        request (flask.Request): HTTP request object.
    Returns:
        The response text or any set of values that can be turned into a
        Response object using
        `make_response <http://flask.pocoo.org/docs/1.0/api/#flask.Flask.make_response>`.
    """
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
