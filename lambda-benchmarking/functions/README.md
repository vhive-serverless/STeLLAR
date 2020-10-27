## Functions Manager
This manager deploys benchmarking-oriented serverless functions to various providers.

### Parameters
- `range` (default "0_300"): How many functions to interact with, and what range of IDs to assign to them.
- `action` (default "deploy"): Should the functions be deployed, removed or updated?
- `provider` (default "aws"): What provider should the manager interact with?
- `sizeBytes` (default "0"): The size of the image to deploy together with functional code, in bytes.
- `logLevel` (default info): Client will use this level for logging information.
