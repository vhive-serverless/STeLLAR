#!/usr/bin/env python

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


def plot_cdfs(args):
    def get_experiment_results(inter_arrival_time):
        burstsize_to_latencies = {}

        experiment_dirs = []
        for dir_path, dir_names, filenames in os.walk(args.path):
            if not dir_names and dir_path.split('IAT')[1].split('s')[0] == inter_arrival_time:
                experiment_dirs.append(dir_path)

        latencies = np.empty((0, 1), float)
        for experiment in experiment_dirs:
            experiment_name = experiment.split('/')[-1]
            burst_size = experiment_name.split('burst')[1].split('-')[0]

            with open(experiment + "/latencies.csv") as file:
                data = pd.read_csv(file)
                read_latencies = data['Client Latency (ms)'].to_numpy()
                sorted_latencies = np.sort(read_latencies)
                burstsize_to_latencies[burst_size] = sorted_latencies

        return burstsize_to_latencies

    def add_composing_cdf(subplot, inter_arrival_time, xlim):
        subplot.set_title(f'{"Warm" if int(inter_arrival_time)<600 else "Cold"} (IAT {inter_arrival_time}s)')
        subplot.set_xlabel('Latency (ms)')
        subplot.set_ylabel('Portion of requests')
        subplot.grid(True)

        subplot.set_xlim([0, xlim])

        burst_sizes = get_experiment_results(inter_arrival_time)

        for size in sorted(burst_sizes):
            quantile = np.arange(len(burst_sizes[size])) / float(len(burst_sizes[size]) - 1)
            subplot.plot(burst_sizes[size], quantile, '--o', markersize=1, label=f'Burst Size {size}')

    title = f'{args.provider} Tail Latency Analysis'
    fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(10, 5))
    fig.suptitle(title)

    add_composing_cdf(axes[0], '3', 800)
    add_composing_cdf(axes[1], '600', 1200)

    plt.legend(loc='lower right')
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig(f'{args.path}/{title}.png')
    fig.savefig(f'{args.path}/{title}.pdf')
    plt.close()

    print("Completed successfully.")
