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

import numpy as np
import pandas as pd
from matplotlib import pyplot as plt


def plot_cdfs(provider, service_time_ms, allocated_memory):
    def add_experiment_subplot(subplot, path):
        with open(path + "/latencies.csv") as file:
            data = pd.read_csv(file)
            read_latencies = data['Client Latency (ms)'].to_numpy()
            sorted_latencies = np.sort(read_latencies)

            quantile = np.arange(len(sorted_latencies)) / float(len(sorted_latencies) - 1)

            subplot.set_xlim([0, 3000.])
            subplot.plot(sorted_latencies, quantile, '--o', markersize=1)

    title = '{0} Cold Starts CDFs (Service Time {1}ms, Memory {2}MB)'.format(provider, service_time_ms,
                                                                             allocated_memory)
    fig = plt.figure(figsize=(12, 5))
    fig.suptitle(title)

    burst_sizes = [1, 100, 200, 500]
    image_sizes_mb = [2.9, 60, 120, 180, 240]
    cols = ['Burst Size {}'.format(col) for col in burst_sizes]
    rows = ['Img. {}MB'.format(row) for row in image_sizes_mb]
    fig, axes = plt.subplots(nrows=len(rows), ncols=len(cols),
                             sharex=True, sharey=True, figsize=(20, 14))
    fig.suptitle(title, fontsize=18)

    for ax, col in zip(axes[0], cols):
        ax.set_title(col, fontsize=16)

    for ax, row in zip(axes[:, -1], rows):
        ax.set_ylabel(row, labelpad=-315, rotation=-90, fontsize=16)

    for ax, col in zip(axes[-1], cols):
        ax.set_xlabel('Latency (ms)')

    for ax, row in zip(axes[:, 0], rows):
        ax.set_ylabel('Fraction')

    path_prefix = 'providers/{0}/{1}MB/st{2}ms'.format(provider, allocated_memory, service_time_ms)
    for col, burst_size in enumerate(burst_sizes):
        for row, image_size_mb in enumerate(image_sizes_mb):
            add_experiment_subplot(axes[row][col],
                                   '{0}/size{1}-img{2}mb'.format(path_prefix, burst_size, image_size_mb))

    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig('{0}/{1}.png'.format(path_prefix, title))
    plt.close()


# First run both experiments `imgsize-burstsize-cold-starts`
# Before running this plotting utility, place your results under, e.g., for
# AWS, providers/AWS/128MB/'st0sec' and 'st1sec'.
plot_cdfs(provider='AWS', service_time_ms='0', allocated_memory='128')
