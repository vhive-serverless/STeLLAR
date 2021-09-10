#!/usr/bin/env python

# MIT License
#
# Copyright (c) 2021 Theodor Amariucai and EASE Lab
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
from matplotlib.lines import Line2D
from matplotlib.ticker import ScalarFormatter


def add_subplot(args, subtitle_percentile, ylabel, subplot, latencies, experiment_type, use_seconds=False):
    def change_to_seconds():
        subplot.set_ylabel('Latency (seconds)')
        return [latency / 1000.0 if latency != np.nan else np.nan for latency in used_latencies]

    if type(latencies['rtt'][experiment_type]) is dict:
        assert latencies['rtt'][experiment_type].keys() == latencies['timestamp_diff'][experiment_type].keys()
    else:
        assert len(latencies['rtt'][experiment_type]) == len(latencies['timestamp_diff'][experiment_type])

    subplot.set_title(subtitle_percentile)
    subplot.set_xlabel('Transfer Size (KB)')
    subplot.set_ylabel(ylabel)

    subplot.set_xscale('log')  # better transfer latencies and bandwidth visualization
    subplot.xaxis.set_major_formatter(ScalarFormatter())
    subplot.grid(True)

    # subplot.set_xticks([int(size/1024.) for size in args.transfer_sizes])

    used_transfer_sizes = args.transfer_sizes
    if max(used_transfer_sizes) >= 1e4:
        subplot.set_xlabel('Transfer Size (MB)')
        used_transfer_sizes = [size / 1024.0 for size in used_transfer_sizes]
    elif max(used_transfer_sizes) >= 1e6:
        subplot.set_xlabel('Transfer Size (GB)')
        used_transfer_sizes = [size / 1024.0 / 1024.0 for size in used_transfer_sizes]

    colors_rtt = []
    for value in sorted(latencies['rtt'][experiment_type]):
        diff = len(used_transfer_sizes) - len(latencies['rtt'][experiment_type][value])
        while diff > 0:
            diff -= 1
            latencies['rtt'][experiment_type][value].append(np.nan)

        used_latencies = latencies['rtt'][experiment_type][value]
        if use_seconds:
            used_latencies = change_to_seconds()

        label = f"{value / 1024.0}GB memory" if experiment_type == 'memory' else f"Chain length {value}"
        rtts_plotted = subplot.plot(used_transfer_sizes, used_latencies, 'o--', label=label)
        colors_rtt.append(rtts_plotted[0].get_color())

        for j, txt in enumerate(used_latencies):
            if not np.isnan(used_latencies[j]):
                subplot.annotate(int(txt), (used_transfer_sizes[j], used_latencies[j]))

    for i, value in enumerate(sorted(latencies['timestamp_diff'][experiment_type])):
        diff = len(used_transfer_sizes) - len(latencies['timestamp_diff'][experiment_type][value])
        while diff > 0:
            diff -= 1
            latencies['timestamp_diff'][experiment_type][value].append(np.nan)

        used_latencies = latencies['timestamp_diff'][experiment_type][value]
        if use_seconds:
            used_latencies = change_to_seconds()

        subplot.plot(used_transfer_sizes, used_latencies, 'o-', color=colors_rtt[i])

        for j, txt in enumerate(used_latencies):
            if not np.isnan(used_latencies[j]):
                subplot.annotate(int(txt), (used_transfer_sizes[j], used_latencies[j]))

    handles, labels = subplot.get_legend_handles_labels()

    if args.provider == "vHive":
        handles, labels = [], []
        legend_color = colors_rtt[0]
    else:
        legend_color = 'black'

    labels.append('Round Trip Time')
    handles.append(Line2D([0], [0], color=legend_color, linewidth=2, linestyle='dotted'))

    labels.append('Internal Timestamp')
    handles.append(Line2D([0], [0], color=legend_color, linewidth=2))

    if "Median" not in subtitle_percentile:
        subplot.legend(handles=handles, labels=labels)


def generate_latency_bandwidth_figures(args, iat_warm_threshold, warm_plots, experiment_type='memory'):
    def load_rtt_and_stampdiff_latencies():
        def fetch_experiment_directories():
            experiment_dirs = []
            for dir_path, dir_names, filenames in os.walk(args.path):
                if not dir_names:  # no subdirectories
                    experiment_dirs.append(dir_path)
            return experiment_dirs

        def read_latencies_median_and_tail():
            with open(experiment + "/latencies.csv") as rtt_file:
                data = pd.read_csv(rtt_file)
                transfer_latencies = data['Client Latency (ms)'].to_numpy()
                sorted_latencies = np.sort(transfer_latencies)

                return sorted_latencies[int(len(sorted_latencies) * 0.5)], sorted_latencies[
                    int(len(sorted_latencies) * args.desired_percentile / 100.0)]

        def read_data_transfer_timestamp_diffs():
            with open(experiment + "/data-transfers.csv") as stamp_file:
                data = pd.read_csv(stamp_file)
                timestamp_start = data['Function 0 Timestamp'].to_numpy()
                timestamp_end = data[f'Function {chain_length - 1} Timestamp'].to_numpy()
                transfer_latencies = timestamp_end - timestamp_start
                sorted_latencies = np.sort(transfer_latencies)

                return sorted_latencies[int(len(sorted_latencies) * 0.5)], sorted_latencies[
                    int(len(sorted_latencies) * args.desired_percentile / 100.0)]

        transfer_sizes_kb = []

        experiment_dirs = fetch_experiment_directories()

        # sort by image size
        experiment_dirs.sort(key=lambda s: float(s.split('payload')[-1].split('KB')[0]))
        # filter by IAT threshold
        if warm_plots:
            experiment_dirs = filter(lambda s: float(s.split('IAT')[1].split('s-')[0]) <= iat_warm_threshold,
                                     experiment_dirs)
        else:
            experiment_dirs = filter(lambda s: float(s.split('IAT')[1].split('s-')[0]) > iat_warm_threshold,
                                     experiment_dirs)
        # filter by payload size
        experiment_dirs = filter(lambda s: float(s.split('payload')[-1].split('KB')[0]) > 0.0, experiment_dirs)

        for experiment in experiment_dirs:
            transfer_sizes_kb.append(float(experiment.split('payload')[-1].split('KB')[0]))
            memory_size = int(experiment.split('memory')[1].split('MB')[0])
            chain_length = int(experiment.split('/')[-1].split('chain')[0])

            median_value, tail_value = read_latencies_median_and_tail()
            if experiment_type == 'chain':
                assign_dictionary_values('rtt', median_value, tail_value, chain_length)
            else:
                assign_dictionary_values('rtt', median_value, tail_value, memory_size)

            median_value, tail_value = read_data_transfer_timestamp_diffs()
            if experiment_type == 'chain':
                assign_dictionary_values('timestamp_diff', median_value, tail_value, chain_length)
            else:
                assign_dictionary_values('timestamp_diff', median_value, tail_value, memory_size)

        args.transfer_sizes = transfer_sizes_kb

    def assign_dictionary_values(latency_type, median_value, tail_value, value):
        if value in median[latency_type][experiment_type]:
            median[latency_type][experiment_type][value].append(median_value)
        else:
            median[latency_type][experiment_type][value] = [median_value]

        if value in percentiles[latency_type][experiment_type]:
            percentiles[latency_type][experiment_type][value].append(tail_value)
        else:
            percentiles[latency_type][experiment_type][value] = [tail_value]

    def generate_transfer_bandwidth_figure():
        title = f'{args.provider} {"Storage" if "storage" in args.path else "Inline"} Transfer Bandwidth'

        fig, axes = plt.subplots(nrows=1, ncols=1, sharey=True, figsize=(7, 5))
        fig.suptitle(title, fontsize=16)

        assert median['rtt'][experiment_type].keys() == median['timestamp_diff'][experiment_type].keys()

        for memory_kb in median['rtt'][experiment_type]:
            median['rtt'][experiment_type][memory_kb] = [(x / 1024) / (y / 1000) for x, y in
                                                         zip(args.transfer_sizes,
                                                             median['rtt'][experiment_type][memory_kb])]

        for memory_kb in median['timestamp_diff'][experiment_type]:
            median['timestamp_diff'][experiment_type][memory_kb] = [(x / 1024) / (y / 1000) for x, y in
                                                                    zip(args.transfer_sizes,
                                                                        median['timestamp_diff'][experiment_type][
                                                                            memory_kb])]

        add_subplot(args, "", 'Network Bandwidth (MB/s)', axes, latencies=median, experiment_type=experiment_type)
        fig.tight_layout(rect=[0, 0, 1, 0.95])
        fig.savefig(f'{args.path}/{title}.png')
        fig.savefig(f'{args.path}/{title}.pdf')
        plt.close()

    def generate_transfer_latency_figure():
        title = f'{args.provider} {"Storage" if "storage" in args.path else "Inline"} Transfer Latency'
        fig, axes = plt.subplots(nrows=1, ncols=1 if args.just_median else 2, sharey=True, figsize=(10, 5))
        fig.suptitle(title, fontsize=16)
        plt.grid(True)

        if args.just_median:
            add_subplot(args, 'Median (50% percentile)', 'Latency (ms)', axes, latencies=median,
                        experiment_type=experiment_type, use_seconds=args.seconds)
        else:
            add_subplot(args, f"{args.desired_percentile}% percentile", 'Latency (ms)', axes[0], latencies=percentiles,
                        experiment_type=experiment_type, use_seconds=args.seconds)
            add_subplot(args, 'Median (50% percentile)', 'Latency (ms)', axes[1], latencies=median,
                        experiment_type=experiment_type, use_seconds=args.seconds)

        fig.tight_layout(rect=[0, 0, 1, 0.95])
        fig.savefig(f'{args.path}/{title}.png')
        fig.savefig(f'{args.path}/{title}.pdf')
        plt.close()

    # dicts from function memories to latencies according to the transfer_size array
    median, percentiles = {}, {}
    median['rtt'], percentiles['rtt'] = {}, {}
    median['timestamp_diff'], percentiles['timestamp_diff'] = {}, {}
    median['rtt']['memory'], percentiles['rtt']['memory'] = {}, {}
    median['rtt']['chain'], percentiles['rtt']['chain'] = {}, {}
    median['timestamp_diff']['memory'], percentiles['timestamp_diff']['memory'] = {}, {}
    median['timestamp_diff']['chain'], percentiles['timestamp_diff']['chain'] = {}, {}
    load_rtt_and_stampdiff_latencies()

    if experiment_type == 'memory':
        args.transfer_sizes = list(dict.fromkeys(args.transfer_sizes))  # remove duplicates

    generate_transfer_latency_figure()
    generate_transfer_bandwidth_figure()


def plot_data_transfer_stats(args):
    args.just_median = False  # True if "vHive" in args.provider else False
    args.desired_percentile = 99 if "vHive" in args.provider else 99

    generate_latency_bandwidth_figures(args=args, iat_warm_threshold=50, warm_plots=True, experiment_type='chain')
