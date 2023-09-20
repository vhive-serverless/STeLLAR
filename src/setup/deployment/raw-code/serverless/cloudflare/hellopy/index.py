from datetime import datetime
def handleRequest(request):
    incr_limit = 0

    if 'queryStringParameters' in request and 'IncrementLimit' in request['queryStringParameters']:
        incr_limit = int(request['queryStringParameters'].get('IncrementLimit', 0))
    elif 'body' in request and JSON.parse(request['body'])['IncrementLimit']:
        incr_limit = int(JSON.parse(request['body'])['IncrementLimit'])

    simulate_work(incr_limit)

    response = JSON.stringify({
        "RequestID": "cloudflare-does-not-specify",
        "TimestampChain": [str(datetime.now())]
    })


    return __new__(Response(response, {
        'headers' : { 'content-type' : 'application/json' }
    }))

def simulate_work(increment):
    # MAXNUM = 6103705
    num = 0
    while num < increment:
        num += 1

addEventListener('fetch', (lambda event: event.respondWith(handleRequest(event.request))))
