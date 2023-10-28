const express = require('express');
const app = express();

const simulateWork = (incrementLimit) => {
  for (let i = 0; i < incrementLimit; i++){}
}

app.get('/', (req, res) => {
  let incrementLimit = 0
  if (req.query.IncrementLimit) {
    incrementLimit = req.query.IncrementLimit
  }
  simulateWork(incrementLimit)

  res.json({
    RequestID: "google-does-not-specify",
    TimestampChain: [Date.now().toString()]
  });
});

const port = process.env.PORT || 8080;

app.use(express.json())
app.listen(port, () => {
  console.log('Function listening on port', port);
});
