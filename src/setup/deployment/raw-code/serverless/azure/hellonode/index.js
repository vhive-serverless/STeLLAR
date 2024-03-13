async function handler(context, request) {
  let q = request.query;
  let incrementLimit = 0;
  if (q.incrementLimit) {
    incrementLimit = parseInt(q.incrementLimit);
  }

  simulateWork(incrementLimit);

  context.res = {
    status: 200,
    body: {
      RequestID: context.invocationId,
      TimestampChain: [Date.now().toString()],
    }
  };
};

const simulateWork = (incrementLimit) => {
  for (let i = 0; i < incrementLimit; i++) { }
};

module.exports = handler
