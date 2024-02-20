import json
import os
import time
import random

import azure.functions as func


def main(req: func.HttpRequest, context: func.Context) -> func.HttpResponse:
    incr_limit = int(req.params.get('IncrementLimit')) if req.params.get('IncrementLimit') else None
    if not incr_limit:
        try:
            req_body = req.get_json()
        except ValueError:
            incr_limit = 0
            pass
        else:
            incr_limit = int(req_body.get('IncrementLimit')) if req_body.get('IncrementLimit') else 0
    else:
        incr_limit = 0

    simulate_work(incr_limit)
    read_filler_file(f"{context.function_directory}/../filler.file")

    return func.HttpResponse(
        body=json.dumps({
            "RequestID": context.invocation_id,
            "TimestampChain": [str(time.time_ns())]
        }, indent=4),
        status_code=200,
        headers={
            "Content-Type": "application/json"
        }
    )


def simulate_work(increment):
    # MAXNUM = 6103705
    num = 0
    while num < increment:
        num += 1


def read_filler_file(path: str) -> None:
    with open(path, 'rb') as f:
        for i in range(1024):
            f.seek(i * 2048)
            f.read(1)
