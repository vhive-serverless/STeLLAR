## Pre-requisites
1. Fetch the deployment branch of STeLLAR, containing the binary as well as other useful configuration files.

```git clone --single-branch --branch deployment https://github.com/ease-lab/STeLLAR.git```

2. Perform some basic update operations, as well as install useful tools (e.g., tmux).
```
cd STeLLAR/scripts/linux && bash setup.sh
```

## Deployment

1. In Azure Functions, functions have to be deployed in separate Function Apps for cold-start measurements to work. Hypothesis: this is because as soon as a function in a Function App is invoked, all functions in that Function App are brought back into memory.