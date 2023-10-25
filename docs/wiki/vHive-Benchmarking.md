## Pre-requisites
- Fetch the deployment branch
```
git clone --single-branch --branch deployment https://github.com/ease-lab/vhive-bench.git
bash setup.sh # do this in directory vhive-bench/scripts/linux
```

- Deploy two _producer-consumer_ functions to vHive:
```
bash deploy_functions.sh # do this in directory vhive-bench/scripts/linux/vhive/burstiness
```

## Benchmarking:
### vHive inline data transfers:
```
sudo ./main -o latency-samples -g endpoints/vhive -c experiments/data-transfer/inline/vhive/warm.json
```

### vHive storage data transfers:

* Create K8s minio storage & bucket:
```
bash minio_setup.sh # do this in directory vhive-bench/scripts/linux/vhive
```

* Run the tool:
```
sudo ./main -o latency-samples -g endpoints -c experiments/data-transfer/storage/vhive-minio/quick-warm-IAT10s.json
```