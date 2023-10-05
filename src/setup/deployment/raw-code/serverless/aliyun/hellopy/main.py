import json
import time
import urllib.parse
from io import BytesIO
from typing import Dict, Any, List


def main(environ, start_response):
    increment_limit = get_increment_limit(environ)
    simulate_work(increment_limit)

    status = "200 OK"
    response_headers = [("Content-type", "application/json")]
    context = environ.get("fc.context")
    response_body = json.dumps({
        "Region": context.region,
        "RequestID": context.request_id,
        "TimestampChain": [str(time.time_ns())],
    })
    start_response(status, response_headers)
    return [response_body]


def get_increment_limit(environ) -> int:
    increment_limit = 0

    query_parameters: Dict[str, List[str]] = urllib.parse.parse_qs(environ.get("QUERY_STRING"))
    if "IncrementLimit" in query_parameters:
        increment_limit = int(query_parameters.get("IncrementLimit")[0])

    content_length = int(environ.get("CONTENT_LENGTH", 0))
    if content_length != 0:
        request_body_buffer: BytesIO = environ.get("wsgi.input")
        request_body_bytes: bytes = request_body_buffer.read(content_length)
        request_body: Dict[str, Any] = json.loads(request_body_bytes)
        if "IncrementLimit" in request_body:
            increment_limit = request_body.get("IncrementLimit")

    return increment_limit


def simulate_work(increment_limit: int) -> None:
    num = 0
    while num < increment_limit:
        num += 1
