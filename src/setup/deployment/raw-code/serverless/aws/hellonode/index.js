// Handler
exports.handler = async function (event, context) {
  let incrementLimit = 0;
  if (event.queryStringParameters.incrementLimit) {
    incrementLimit = event.queryStringParameters.incrementLimit;
  }
  simulateWork(incrementLimit);
  const res = {
    statusCode: 200,
    headers: { "Content-Type": "application/json" },
    body: {
      RequestID: context.aws_request_id,
      TimestampChain: [Date.now().toString()],
    },
  };

  return JSON.stringify(res);
};

const simulateWork = (incrementLimit) => {
  for (let i = 0; i < incrementLimit; i++) {}
};
