# MIT License
#
# Copyright (c) 2020 Theodor Amariucai
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


def plot_cold_starts(provider, service_time, allocated_memory):
    image_size_mb = [2.9, 60, 120, 180, 240]

    def add_percentile_subplot(subtitle_percentile, subplot, bs_1, bs_100, bs_200, bs_500):
        subplot.set_title(subtitle_percentile)
        subplot.set_xlabel('Image Size (MB)')
        subplot.set_ylabel('Latency (ms)')

        subplot.plot(image_size_mb, bs_1, 'o-', label='burst_size_1')
        subplot.plot(image_size_mb, bs_100, 'o-', label='burst_size_100')
        subplot.plot(image_size_mb, bs_200, 'o-', label='burst_size_200')
        subplot.plot(image_size_mb, bs_500, 'o-', label='burst_size_500')
        subplot.legend(loc='upper left')

    def load_experiment_results(path):
        experiment_dirs = []
        for dir_path, dir_names, filenames in os.walk(path):
            if not dir_names:  # no subdirectories
                experiment_dirs.append(dir_path)

        # sort by image size
        experiment_dirs.sort(key=lambda s: float(s.split('img')[1].split('mb')[0]))

        for experiment in experiment_dirs:
            experiment_name = experiment.split('/')[4]
            burst_size = int(experiment_name.split('-')[0].split('size')[1])
            # image_size = int(experiment_name.split('-')[1].split('img')[1].split('mb')[0])
            with open(experiment + "/latencies.csv") as file:
                data = pd.read_csv(file)
                read_latencies = data['Client Latency (ms)'].to_numpy()
                sorted_latencies = np.sort(read_latencies)

                median[burst_size].append(sorted_latencies[int(len(sorted_latencies) * 0.5)])  # 50%ile = median
                percentiles[burst_size].append(sorted_latencies[int(len(sorted_latencies) * 0.95)])  # 95%ile

    # dicts from batch sizes to latencies according to the image_size_mb array
    percentiles = {
        1: [],
        100: [],
        200: [],
        500: [],
    }
    median = {
        1: [],
        100: [],
        200: [],
        500: [],
    }

    path_prefix = 'providers/{0}/{1}MB/st{2}'.format(provider, allocated_memory, service_time)
    load_experiment_results(path_prefix)

    title = '{0} Cold Starts (Service Time {1}, Memory {2}MB)'.format(provider, service_time, allocated_memory)

    fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(12, 5))
    fig.suptitle(title)

    add_percentile_subplot('95% percentile', axes[0], percentiles[1], percentiles[100], percentiles[200],
                           percentiles[500])
    add_percentile_subplot('Median (50% percentile)', axes[1], median[1], median[100], median[200], median[500])

    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig('{0}/{1}.png'.format(path_prefix, title))
    plt.close()


# First run both experiments `imgsize-burstsize-cold-starts`
# Before running this plotting utility, place your results under, e.g., for
# AWS, providers/AWS/128MB/'st0sec' and 'st1sec'.
plot_cold_starts(provider='AWS', service_time='1s', allocated_memory='1536')
