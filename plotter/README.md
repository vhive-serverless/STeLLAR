# vHive-bench plotting utility
These auxiliary Python scripts are used to visualize the benchmarking results.

## Usage examples:
1. CPU stats

`python plot.py --type cpustats --path ../latency-samples/cloudlab/cpu-stats/AWS/Tuesday, 16-Mar-21 05:36:42 MDT/`

2. Transfers:

`python plot.py --type transfer --path ../latency-samples/cloudlab/data-transfer/AWS/inline/Saturday, 27-Mar-21 10:33:31 MDT/`

3. Burstiness CDFs

`python plot.py --type burstiness --path ../latency-samples/cloudlab/burstiness/AWS/Tuesday\,\ 16-Mar-21\ 08\:19\:39\ MDT/`

4. Cold starts image size:

`python plot.py --type imgsize --path ../latency-samples/cloudlab/image-size/AWS/1536MB\ memory\,\ st1s\ imgsize\ experiment/`
