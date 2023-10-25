## Pre-requisites
1. [Create an AWS account](https://portal.aws.amazon.com/billing/signup#/start), then go to `My security credentials` -> `Access keys (access key ID and secret access key)` -> `Create New Access Key`. Take note of your Access Key ID and Secret Access Key.
2. Fetch the deployment branch of STeLLAR, containing the binary as well as other useful configuration files.

```git clone --single-branch --branch deployment https://github.com/ease-lab/STeLLAR.git```

3. Perform some basic update operations, as well as install useful tools (tmux).
```
cd STeLLAR/scripts && bash setup.sh
cd aws && bash setup.sh # Note: here please use your AWS Access Key created at step 1
```

## Setup
1. **Add permissions for the Lambda functions:** IAM console -> Roles -> Create Role -> Lambda -> Select policies `AWSLambdaFullAccess`, `AWSLambdaBasicExecutionRole` (for logging), `AWSLambdaRole` (for triggering other lambdas in a producer-consumer chain scenario) and `AmazonS3FullAccess ` (for inter-function storage transfers using S3) -> Enter name **exactly** as `LambdaProducerConsumer`. Take note of your 12-digit user ARN number: `arn:aws:iam::12-DIGIT-ARN-NUMBER:role/LambdaProducerConsumer`).
3. **Add permissions for the STeLLAR client:** IAM console -> Users -> Create User. Attach `AWSLambda_FullAccess` and `AmazonS3FullAccess` to this user, and then create and attach another policy, e.g., named `APIGatewayFull`, which has `API Gateway; Full Access; All resources` and `API Gateway V2; Full Access; All resources`.

~3. **Local client authentication:** IAM console -> Users -> Summary -> Security credentials -> Create access key. **Make sure to keep creating keys until there are no forward slashes in them**. Take note of your Access key ID (e.g., `AKIAU4EZSQEM3S42BXY5`) and Secret access key (e.g., `vAbqZTA3MpxctAi4zvS5QW4Qvvhpkg53lALgDUDV`). Configure your local AWS CLI by running `aws configure`, use your Access Key ID and Secret access key. The default region name used in the client is `us-west-1`.~

### To further enable function [deployments via container images](https://docs.aws.amazon.com/lambda/latest/dg/images-create.html#images-create-1):
1. Create a private Amazon ECR repository with the name `vhive-bench`.
2. IAM console -> Users. Select your user created in step 2 and attach the `AmazonEC2ContainerRegistryFullAccess` permission policy to it.

## Benchmarking:
We have provided ready-made configurations that helped us generate our [published results](https://github.com/ease-lab/STeLLAR/blob/main/docs/STeLLAR_IISWC21.pdf), but you can easily [design your own experiments](https://github.com/ease-lab/STeLLAR/wiki/Customize-Experiments).

### Inline data transfers:
```
sudo ./main -a 12-DIGIT-ARN-NUMBER -o latency-samples -g endpoints -c experiments/data-transfer/inline/aws/quick-warm-IAT10s.json
```

### Storage data transfers:
```
sudo ./main -a 12-DIGIT-ARN-NUMBER -o latency-samples -g endpoints -c experiments/data-transfer/storage/
```