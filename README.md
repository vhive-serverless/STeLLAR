# vHive-bench ![CI Result](https://github.com/ease-lab/vhive-bench/workflows/Go/badge.svg?branch=master)
A framework for benchmarking the performance of popular serverless platforms. 

## Design
![design](design/diagram.png)

## Flow Chart
![flow chart](design/flow-chart.png)

## Common problems

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