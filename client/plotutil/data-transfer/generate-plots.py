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
from matplotlib.ticker import ScalarFormatter


def add_subplot(subtitle_percentile, ylabel, subplot, rtt_latencies, stamp_latencies):
    subplot.set_title(subtitle_percentile)
    subplot.set_xlabel('Transfer Size (KB)')
    subplot.set_ylabel(ylabel)

    subplot.set_xscale('log')  # better transfer latencies and bandwidth visualization
#     subplot.xaxis.set_major_formatter(ScalarFormatter())

    colors_rtt = []
    for memory_mb in rtt_latencies:
        line_label = "RTT {0}GB (memory size)".format(memory_mb/1024.0)

        if provider == "vhive":
            line_label = "RTT"

        rtts_plotted = subplot.plot(transfer_size_kb, rtt_latencies[memory_mb], 'o-', label=line_label)
        colors_rtt.append(rtts_plotted[0].get_color())

        for i, txt in enumerate(rtt_latencies[memory_mb]):
            if not np.isnan(rtt_latencies[memory_mb][i]):
                subplot.annotate(int(txt), (transfer_size_kb[i], rtt_latencies[memory_mb][i]))

    for i, memory_mb in enumerate(stamp_latencies):
        line_label = "Internal Timestamp {0}GB (memory size)".format(memory_mb/1024.0)

        if provider == "vhive":
            line_label = "Timestamp Delta"

        subplot.plot(transfer_size_kb, stamp_latencies[memory_mb], 'o--', label=line_label, color=colors_rtt[i])

        for i, txt in enumerate(stamp_latencies[memory_mb]):
            if not np.isnan(stamp_latencies[memory_mb][i]):
                subplot.annotate(int(txt), (transfer_size_kb[i], stamp_latencies[memory_mb][i]))

    subplot.legend(loc='upper left')


def load_experiment_results(path, inter_arrival_time):
    # dicts from function memories to latencies according to the transfer_size array
    rtt_median = {}
    rtt_percentiles = {}
    stamp_diff_median = {}
    stamp_diff_percentiles = {}

    experiment_dirs = []
    for dir_path, dir_names, filenames in os.walk(path):
        if not dir_names:  # no subdirectories
            experiment_dirs.append(dir_path)

    # sort by image size
    experiment_dirs.sort(key=lambda s: float(s.split('-')[3].split('KBpayload')[0]))
    experiment_dirs = filter(lambda s: s.split('IAT')[1].split('-')[0] == inter_arrival_time, experiment_dirs)

    for experiment in experiment_dirs:
        memory_size = int(experiment.split('-')[1].split('MB')[0])

        with open(experiment + "/latencies.csv") as rtt_file:
            data = pd.read_csv(rtt_file)
            transfer_latencies = data['Client Latency (ms)'].to_numpy()
            sorted_latencies = np.sort(transfer_latencies)

            median_value = sorted_latencies[int(len(sorted_latencies) * 0.5)]
            if memory_size in rtt_median:
                rtt_median[memory_size].append(median_value)  # 50%ile = median
            else:
                rtt_median[memory_size] = [median_value]

            tail_value = sorted_latencies[int(len(sorted_latencies) * desired_percentile / 100.0)]
            if memory_size in rtt_percentiles:
                rtt_percentiles[memory_size].append(tail_value)  # 50%ile = median
            else:
                rtt_percentiles[memory_size] = [tail_value]

        with open(experiment + "/data-transfers.csv") as stamp_file:
            data = pd.read_csv(stamp_file)
            timestamp1 = data['Function 0 Timestamp'].to_numpy()
            timestamp2 = data['Function 1 Timestamp'].to_numpy()
            transfer_latencies = timestamp2 - timestamp1
            sorted_latencies = np.sort(transfer_latencies)

            median_value = sorted_latencies[int(len(sorted_latencies) * 0.5)]
            if memory_size in stamp_diff_median:
                stamp_diff_median[memory_size].append(median_value)  # 50%ile = median
            else:
                stamp_diff_median[memory_size] = [median_value]

            tail_value = sorted_latencies[int(len(sorted_latencies) * desired_percentile / 100.0)]
            if memory_size in stamp_diff_percentiles:
                stamp_diff_percentiles[memory_size].append(tail_value)  # 50%ile = median
            else:
                stamp_diff_percentiles[memory_size] = [tail_value]

    return rtt_median, rtt_percentiles, stamp_diff_median, stamp_diff_percentiles


def generate_transfer_bandwidth_figure(path, inter_arrival_time, rtt_median, stamp_diff_median):
    if("storage" in path):
        title = '{0} Storage Bandwidth (IAT {1})'.format(provider, inter_arrival_time)
    else:
        title = '{0} Inline Bandwidth (IAT {1})'.format(provider, inter_arrival_time)

    fig, axes = plt.subplots(nrows=1, ncols=1, sharey=True, figsize=(8, 5))
    fig.suptitle(title)

    for memory_kb in rtt_median:
        rtt_median[memory_kb] = [x / y * 1000 / 1024 for x, y in zip(transfer_size_kb, rtt_median[memory_kb])]

    for memory_kb in stamp_diff_median:
        stamp_diff_median[memory_kb] = [x / y * 1000 / 1024 for x, y in zip(transfer_size_kb, stamp_diff_median[memory_kb])]

    add_subplot("", 'Network Bandwidth (MB/s)', axes, rtt_median, stamp_diff_median)
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig('{0}/{1}.png'.format(path, title))
    plt.close()


def generate_transfer_latency_figure(path, inter_arrival_time, rtt_median, rtt_percentiles, stamp_diff_median, stamp_diff_percentiles):
    if("storage" in path):
        title = '{0} Storage Transfer (IAT {1})'.format(provider, inter_arrival_time)
    else:
        title = '{0} Inline Transfer (IAT {1})'.format(provider, inter_arrival_time)

    if just_median:
        fig, axes = plt.subplots(nrows=1, ncols=1, sharey=True, figsize=(12, 5))
        fig.suptitle(title)
        add_subplot('Median (50% percentile)', 'Transfer Latency (ms)', axes, rtt_median, stamp_diff_median)
    else:
        fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(12, 5))
        fig.suptitle(title)
        add_subplot("{0}% percentile".format(desired_percentile), 'Transfer Latency (ms)', axes[0], rtt_percentiles, stamp_diff_percentiles)
        add_subplot('Median (50% percentile)', 'Transfer Latency (ms)', axes[1], rtt_median, stamp_diff_median)

    fig.savefig('{0}/{1}.png'.format(path, title))
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    plt.close()


def generate_figures(inter_arrival_time):
    rtt_median, rtt_percentiles, stamp_diff_median, stamp_diff_percentiles = load_experiment_results(path_prefix, inter_arrival_time)

    generate_transfer_latency_figure(path_prefix, inter_arrival_time, rtt_median, rtt_percentiles, stamp_diff_median, stamp_diff_percentiles)
    generate_transfer_bandwidth_figure(path_prefix, inter_arrival_time, rtt_median, stamp_diff_median)


provider = 'AWS'
just_median = True
# First run both experiments `data-transfer/cold-IAT600s-reduced-samples.json` and
# `data-transfer/warm-IAT10s.json`. Then remove the 0KB results (log scale), and
# place the rest of them under, e.g., for AWS inline, `providers/AWS/inline`
desired_percentile = 99
transfer_size_kb = [64000, 128000, 256000, 512000, 1048576]
# [1, 10, 100, 1024, 10240, 102400]
# [1, 10, 100, 200, 400, 600, 800, 1024, 2048, 3072, 4096]
path_prefix = 'providers/{0}/storage_large'.format(provider)
# generate_figures('600s')
generate_figures('10s')