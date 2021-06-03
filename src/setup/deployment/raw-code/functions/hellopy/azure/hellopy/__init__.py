import logging

import json
import azure.functions as func


def main(req: func.HttpRequest, context: func.Context) -> func.HttpResponse:
    # name = req.params.get('name')
    # if not name:
    #     try:
    #         req_body = req.get_json()
    #     except ValueError:
    #         pass
    #     else:
    #         name = req_body.get('name')

    return func.HttpResponse(
        status_code=200,
        headers={
            "Content-Type": "application/json"
        },
        body=json.dumps({
            # "Region ": json_region,
            "RequestID": context.invocation_id,
            "TimestampChain": ['0'],
        })
    )
