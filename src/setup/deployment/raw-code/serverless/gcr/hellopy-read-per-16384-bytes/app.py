import json
import os
import time
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
    with open(path, 'rb') as f:
        for i in range(1024):
            f.seek(i * 16384)
            f.read(1)


if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=int(os.environ.get('PORT', 8080)))
