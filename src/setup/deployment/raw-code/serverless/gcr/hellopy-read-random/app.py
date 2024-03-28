import json
import os
import time
import random
from flask import Flask, request

app = Flask(__name__)


@app.route('/')
def hello_world():
    incr_limit = 0
    if request.args and 'incrementLimit' in request.args:
        incr_limit = request.args.get('incrementLimit')

    read_filler_file("./filler.file")
    simulate_work(incr_limit)

    response = {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": {
            "RequestID": "gcr-does-not-specify",
            "TimestampChain": [str(time.time_ns())],
        }
    }

    return json.dumps(response, indent=4)


def simulate_work(incr):
    num = 0
    while num < incr:
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


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
