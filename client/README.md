## Benchmarking Client ![Go](https://github.com/ease-lab/vhive-bench/workflows/Go/badge.svg?branch=master)
This client tests the performance of 
AWS Lambda busy-spinning microVM functions by sending requests and benchmarking the
latencies.

## Design
![design](design/diagram.png)

## Flow Chart
![design](design/flow-chart.png)

## Provider Limitations

### AWS
- Code storage limit
```
Cannot update function code: CodeStorageExceededException: Code storage limit exceeded.
{
  RespMetadata: {
    StatusCode: 400,
    RequestID: "886339b1-63ae-4f80-a923-7c1ed4201b6e"
  },
  Message_: "Code storage limit exceeded.",
  Type: "User"
}
```

- Regional APIs limit `600`