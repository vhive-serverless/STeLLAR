require 'json'

def handler(event:, context:)
  incrementLimit = 0
  if event['queryStringParameters'] != nil
    if event['queryStringParameters']['IncrementLimit'] != nil
      incrementLimit = event['queryStringParameters']['IncrementLimit'].to_i
    end
  end

  simulateWork(incrementLimit)

  { RequestID: context.aws_request_id, TimestampChain: [ DateTime.now.strftime('%Q') ] }
end

def simulateWork(incrementLimit)
  i = 0
  while i < incrementLimit
    i += 1
  end
end
