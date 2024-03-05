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

    simulate_work(incr_limit)
    read_filler_file("./filler.file")

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
    page_size_in_bytes = 4096
    total_bytes_to_read = 2400
    sequential_pages_to_read = 3

    number_of_random_positions = total_bytes_to_read // sequential_pages_to_read

    file_size_in_bytes = os.stat(path).st_size
    random_position_upper_limit = file_size_in_bytes - (sequential_pages_to_read * page_size_in_bytes)

    with open(path, 'rb') as f:
        for _ in range(number_of_random_positions):
            random_position = random.randrange(0, random_position_upper_limit)
            for i in range(sequential_pages_to_read):
                f.seek(random_position + (i * page_size_in_bytes))
                f.read(1)


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
