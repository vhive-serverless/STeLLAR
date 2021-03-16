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
import copy
import os

import numpy as np
import pandas as pd
from matplotlib import pyplot as plt


def load_experiment_results(args):
    image_sizes_mb = []
    percentiles = {}
    median = {}
    latencies = {}
    quantiles = {}

    experiment_dirs = []
    for dir_path, dir_names, filenames in os.walk(args.path):
        if not dir_names:  # no subdirectories
            experiment_dirs.append(dir_path)

    # sort by image size
    experiment_dirs.sort(key=lambda s: float(s.split('/')[-1].split('img')[1].split('MB')[0]))

    for experiment in experiment_dirs:
        experiment_name = experiment.split('/')[-1]
        burst_size = int(experiment_name.split('burst')[1].split('-')[0])
        image_size = int(experiment_name.split('img')[1].split('MB')[0])
        image_sizes_mb.append(image_size)

        with open(experiment + "/latencies.csv") as file:
            data = pd.read_csv(file)
            read_latencies = data['Client Latency (ms)'].to_numpy()
            sorted_latencies = np.sort(read_latencies)

            median_value = sorted_latencies[int(len(sorted_latencies) * 0.5)]
            if burst_size in median:
                median[burst_size].append(median_value)  # 50%ile = median
            else:
                median[burst_size] = [median_value]

            tail_value = sorted_latencies[int(len(sorted_latencies) * args.desired_percentile / 100.0)]
            if burst_size in percentiles:
                percentiles[burst_size].append(tail_value)
            else:
                percentiles[burst_size] = [tail_value]

            latencies[(burst_size, image_size)] = sorted_latencies
            quantiles[(burst_size, image_size)] = np.arange(len(sorted_latencies)) / float(len(sorted_latencies) - 1)

    return image_sizes_mb, percentiles, median, quantiles, latencies


def plot_imgsize_experiment(args):
    def add_percentile_subplot(subtitle_percentile, subplot, chosen_percentiles, image_sizes_mb):
        subplot.set_title(subtitle_percentile)
        subplot.set_xlabel('Image Size (MB)')
        subplot.set_ylabel('Latency (ms)')
        subplot.grid(True)

        for burst_size in sorted(args.burst_sizes):
            subplot.plot(image_sizes_mb, chosen_percentiles[burst_size], 'o-', label=f"Burst Size {burst_size}")

            for i, txt in enumerate(chosen_percentiles[burst_size]):
                if not np.isnan(chosen_percentiles[burst_size][i]):
                    subplot.annotate(int(txt), (image_sizes_mb[i], chosen_percentiles[burst_size][i]))

        subplot.legend(loc='upper left')

    title = f'{args.provider} Cold Starts - Image Size Experiments (Service Time {args.service_time})'

    fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(10, 5))
    fig.suptitle(title)

    add_percentile_subplot(f"{args.desired_percentile}% percentile", axes[0], args.percentiles, args.image_sizes_mb)
    add_percentile_subplot('Median (50% percentile)', axes[1], args.median, args.image_sizes_mb)

    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig(f'{args.path}/{title}.png')
    fig.savefig(f'{args.path}/{title}.pdf')
    plt.close()


def plot_imgsize_cdfs(args):
    def plot_legend():
        handles, labels = axes[col].get_legend_handles_labels()  # obtain the handles and labels from the figure
        handles = [copy.copy(ha) for ha in handles]  # copy the handles
        [ha.set_linewidth(7) for ha in handles]  # set the linewidths to the copies
        # put the copies into the legend
        axes[col].legend(handles=handles, labels=labels, loc='lower right', prop={'size': 14})

    title = f'{args.provider} Cold Starts CDFs (Service Time {args.service_time})'

    fig = plt.figure(figsize=(12, 5))
    fig.suptitle(title)

    fig, axes = plt.subplots(nrows=1, ncols=len(args.burst_sizes), sharex=True, sharey=True, figsize=(20, 5))
    fig.suptitle(title, fontsize=18)

    for col, burst_size in enumerate(sorted(args.burst_sizes)):
        axes[col].set_title(f'Burst Size {burst_size}', fontsize=15)
        axes[col].set_ylabel('Fraction')
        axes[col].set_xlabel('Latency (ms)')
        axes[col].set_xlim([0, 5000. if args.service_time == '1s' else 3800.])

        for image_size in sorted(args.image_sizes_mb):
            axes[col].plot(args.latencies[(burst_size, image_size)], args.quantiles[(burst_size, image_size)], '--o',
                           markersize=1, label=f"Image Size {image_size}MB")

        if col == len(args.burst_sizes) - 1:
            plot_legend()

    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig(f'{args.path}/{title}.png')
    fig.savefig(f'{args.path}/{title}.pdf')
    plt.close()


def plot_imgsize_stats(args):
    args.desired_percentile = 95

    args.service_time = '1s' if "service-time-1s" in args.path else '0ms'

    args.image_sizes_mb, args.percentiles, args.median, args.quantiles, args.latencies = load_experiment_results(args)
    args.image_sizes_mb = list(dict.fromkeys(args.image_sizes_mb))  # remove duplicates

    args.burst_sizes = []
    for burst_size in sorted(args.percentiles):
        args.burst_sizes.append(burst_size)

    plot_imgsize_experiment(args)
    plot_imgsize_cdfs(args)
