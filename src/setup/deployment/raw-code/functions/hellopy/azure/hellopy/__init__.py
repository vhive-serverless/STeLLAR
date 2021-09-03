import logging

import json
import azure.functions as func


def main(req: func.HttpRequest, context: func.Context) -> func.HttpResponse:
    return func.HttpResponse(
        status_code=200,
        headers={
            "Content-Type": "application/json"
        },
        body=json.dumps({
            "RequestID": context.invocation_id,
            "TimestampChain": ['0'],
        })
    )
