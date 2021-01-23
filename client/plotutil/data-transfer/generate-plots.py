# MIT License
#
# Copyright (c) 2021 Theodor Amariucai
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

import os

import numpy as np
import pandas as pd
from matplotlib import pyplot as plt


def add_subplot(subtitle_percentile, ylabel, subplot, ms_128, ms_1536, ms_10240):
    subplot.set_title(subtitle_percentile)
    subplot.set_xlabel('Transfer Size (KB)')
    subplot.set_ylabel(ylabel)

    subplot.set_xscale('log')  # better transfer latencies and bandwidth visualization

    subplot.plot(transfer_size_kb, ms_128, 'o-', label='128MB (memory size)')
    for i, txt in enumerate(ms_128):
        subplot.annotate(int(txt), (transfer_size_kb[i], ms_128[i]))
    subplot.plot(transfer_size_kb, ms_1536, 'o-', label='1536MB (memory size)')
    for i, txt in enumerate(ms_1536):
        subplot.annotate(int(txt), (transfer_size_kb[i], ms_1536[i]))
    subplot.plot(transfer_size_kb, ms_10240, 'o-', label='10GB (memory size)')
    for i, txt in enumerate(ms_10240):
        subplot.annotate(int(txt), (transfer_size_kb[i], ms_10240[i]))
    subplot.legend(loc='upper left')


def load_experiment_results(path, inter_arrival_time):
    # dicts from function memories to latencies according to the transfer_size array
    percentiles = {
        128: [],
        1536: [],
        10240: [],
    }
    median = {
        128: [],
        1536: [],
        10240: [],
    }

    experiment_dirs = []
    for dir_path, dir_names, filenames in os.walk(path):
        if not dir_names:  # no subdirectories
            experiment_dirs.append(dir_path)

    # sort by image size
    experiment_dirs.sort(key=lambda s: float(s.split('-')[3].split('KBpayload')[0]))
    experiment_dirs = filter(lambda s: s.split('IAT')[1].split('-')[0] == inter_arrival_time, experiment_dirs)

    for experiment in experiment_dirs:
        memory_size = int(experiment.split('-')[1].split('MB')[0])
        with open(experiment + "/data-transfers.csv") as file:
            data = pd.read_csv(file)
            timestamp1 = data['Function 0 Timestamp'].to_numpy()
            timestamp2 = data['Function 1 Timestamp'].to_numpy()
            transfer_latencies = timestamp2 - timestamp1
            sorted_latencies = np.sort(transfer_latencies)

            median[memory_size].append(sorted_latencies[int(len(sorted_latencies) * 0.5)])  # 50%ile = median
            percentiles[memory_size].append(
                sorted_latencies[int(len(sorted_latencies) * desired_percentile / 100.0)])  # 95%ile

    return median, percentiles


def generate_transfer_bandwidth_figure(path, inter_arrival_time, median):
    title = '{0} Direct JSON Transfer Bandwidth (IAT {1})'.format(provider, inter_arrival_time)
    fig, axes = plt.subplots(nrows=1, ncols=1, sharey=True, figsize=(8, 5))
    fig.suptitle(title)
    add_subplot("", 'Network Bandwidth (KB/s)', axes,
                [x / y * 1000 for x, y in zip(transfer_size_kb, median[128])],
                [x / y * 1000 for x, y in zip(transfer_size_kb, median[1536])],
                [x / y * 1000 for x, y in zip(transfer_size_kb, median[10240])])
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig('{0}/{1}.png'.format(path, title))
    plt.close()


def generate_transfer_latency_figure(path, inter_arrival_time, median, percentiles):
    title = '{0} Direct JSON Transfer (IAT {1})'.format(provider, inter_arrival_time)
    fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(12, 5))
    fig.suptitle(title)
    add_subplot("{0}% percentile".format(desired_percentile), 'Transfer Latency (ms)',
                axes[0], percentiles[128], percentiles[1536], percentiles[10240])
    add_subplot('Median (50% percentile)', 'Transfer Latency (ms)',
                axes[1], median[128], median[1536], median[10240])
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig('{0}/{1}.png'.format(path, title))
    plt.close()


def generate_figures(inter_arrival_time):
    median, percentiles = load_experiment_results(path_prefix, inter_arrival_time)
    generate_transfer_latency_figure(path_prefix, inter_arrival_time, median, percentiles)
    generate_transfer_bandwidth_figure(path_prefix, inter_arrival_time, median)


# First run both experiments `data-transfer/cold-IAT600s-reduced-samples.json` and
# `data-transfer/warm-IAT10s.json`. Then remove the 0KB results, and
# place the rest of them under, e.g., for AWS, `providers/AWS/`
provider = 'AWS'
desired_percentile = 99
transfer_size_kb = [1, 10, 100, 1024]
path_prefix = 'providers/{0}'.format(provider)
generate_figures('600s')
generate_figures('10s')
