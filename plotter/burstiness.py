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
import statistics

import numpy as np
import pandas as pd
from matplotlib import pyplot as plt
from matplotlib.lines import Line2D


def plot_cdfs(args):
    def plot_composing_cdf_return_latencies(subplot, iat_interval, xstart, xend):
        desired_percentile = 0.99
        is_warm = int(iat_interval[1]) < 600

        subplot.set_title(f'{"Warm" if is_warm else "Cold"} (IAT {iat_interval}s)')
        subplot.set_xlabel('Latency (ms)')
        subplot.set_ylabel('Portion of requests')
        subplot.grid(True)

        subplot.set_xlim([xstart, xend])

        burst_sizes = get_experiment_results(iat_interval)

        plotting_annotation_index = 1
        for size in sorted(burst_sizes):
            latencies = burst_sizes[size]
            if is_warm or size == '1':  # remove cold latencies from warm instance experiments
                latencies = latencies[:-int(size)]  # remove extra cold latencies

            quantile = np.arange(len(latencies)) / float(len(latencies) - 1)
            recent = subplot.plot(latencies, quantile, '--o', markersize=3, label=f'Burst Size {size}',
                                  markerfacecolor='none')

            median_latency = latencies[int(0.5 * len(latencies))]
            subplot.axvline(x=median_latency, color=recent[-1].get_color(), linestyle='--')
            subplot.annotate(f'{median_latency:0.0f}ms',
                             (min(int(median_latency) * 1.1, int(median_latency) + 2),
                              0.5 + plotting_annotation_index * 0.1),
                             color='black')

            tail_latency = latencies[int(desired_percentile * len(latencies))]
            subplot.axvline(x=tail_latency, color=recent[-1].get_color(), linestyle='--')
            subplot.annotate(f'{tail_latency:0.0f}ms', (
                min(int(tail_latency) * 1.1, int(tail_latency) + 2), 0.1 + plotting_annotation_index * 0.1),
                             color='red')

            plotting_annotation_index += 1

        return burst_sizes

    def plot_dual_cdf(path, latencies_dict, burst_size):
        _fig = plt.figure(figsize=(5, 5))
        _fig.suptitle(f'Burst size {burst_size} ({args.provider})')
        plt.xlabel('Latency (ms)')
        plt.ylabel('Portion of requests')
        plt.grid(True)

        for iat in ['600s', '3s']:
            latencies = latencies_dict[iat][burst_size]
            if iat == '3s' or burst_size == '1':
                latencies = latencies[:-(int(burst_size) + 5)]  # remove extra cold latencies + outliers

            quantile = np.arange(len(latencies)) / float(len(latencies) - 1)
            recent = plt.plot(latencies, quantile, '--o', markersize=4, markerfacecolor='none',
                              label=f'{"Warm" if iat == "3s" else "Cold"} (IAT {iat})')

            print(f'Max latency {latencies[-1]}, stddev {statistics.stdev(latencies)}')

            median_latency = latencies[int(0.5 * len(latencies))]
            plt.axvline(x=median_latency, color=recent[-1].get_color(), linestyle='--')
            plt.annotate(f'{median_latency:0.0f}ms', (int(median_latency) + 2, 0.6 if iat == "3s" else 0.8),
                         color='black')

            tail_latency = latencies[int(0.99 * len(latencies))]
            plt.axvline(x=tail_latency, color=recent[-1].get_color(), linestyle='--')
            plt.annotate(f'{tail_latency:0.0f}ms', (int(tail_latency) + 2, 0.2 if iat == "3s" else 0.4), color='red')

        plt.legend(loc='lower right')
        _fig.savefig(f'{path}/burst{burst_size}-dual-IAT-CDF.png')
        _fig.savefig(f'{path}/burst{burst_size}-dual-IAT-CDF.pdf')
        plt.close()

    def plot_individual_cdf(path, inter_arrival_time, latencies, size):
        desired_percentile = 0.99

        if 'warm' in path or size == '1':  # remove cold latencies from warm instance experiments
            latencies = latencies[:-int(size)]

        _fig = plt.figure(figsize=(5, 5))
        _fig.suptitle(f'Burst size {size}, IAT ~{inter_arrival_time}s ({args.provider})')
        plt.xlabel('Latency (ms)')
        plt.ylabel('Portion of requests')
        plt.grid(True)

        quantile = np.arange(len(latencies)) / float(len(latencies) - 1)
        recent = plt.plot(latencies, quantile, '--o', markersize=4, markerfacecolor='none', color='black')

        median_latency = latencies[int(0.5 * len(latencies))]
        plt.axvline(x=median_latency, color=recent[-1].get_color(), linestyle='--')
        plt.annotate(f'{median_latency:0.0f}ms',
                     (min(int(median_latency) * 1.1, int(median_latency) + 2), 0.5 if 'warm' in path else 0.75),
                     color='black')

        tail_latency = latencies[int(desired_percentile * len(latencies))]
        plt.axvline(x=tail_latency, color='red', linestyle='--')
        plt.annotate(f'{tail_latency:0.0f}ms', (min(int(tail_latency) * 1.1, int(tail_latency) + 2), 0.25),
                     color='red')

        handles, labels = [], []

        labels.append('Average')
        handles.append(Line2D([0], [0], color=recent[-1].get_color(), linewidth=2, linestyle='dotted'))

        labels.append(f'{int(desired_percentile * 100)}%ile')
        handles.append(Line2D([0], [0], color='red', linewidth=2, linestyle='dotted'))

        legend = plt.legend(handles=handles, labels=labels, loc='lower right')
        legend.get_texts()[1].set_color("red")

        _fig.savefig(f'{path}/empirical-CDF.png')
        _fig.savefig(f'{path}/empirical-CDF.pdf')
        plt.close()

    def get_experiment_results(iat_interval):
        burstsize_to_latencies = {}

        experiment_dirs = []
        for dir_path, dir_names, filenames in os.walk(args.path):
            iat = int(dir_path.split('IAT')[1].split('s')[0])
            if not dir_names and iat_interval[0] <= iat <= iat_interval[1]:
                experiment_dirs.append(dir_path)

        for experiment in experiment_dirs:
            experiment_name = experiment.split('/')[-1]
            burst_size = experiment_name.split('burst')[1].split('-')[0]

            with open(experiment + "/latencies.csv") as file:
                data = pd.read_csv(file)

                if args.provider.lower() != "google":
                    data.fillna('', inplace=True)
                    data = data[data["Request ID"].str.len() > 0]
                    print(f'Experiment "{experiment}" had {len(data)} samples not missing/timed out/404!')

                if args.provider.lower() == "azure":
                    if iat_interval[0] == 600:
                        data = data[data['Client Latency (ms)'] > 200]  # filter warm reqs from cold reqs

                if args.provider.lower() == "google":
                    if iat_interval[0] == 600:
                        data = data[data['Client Latency (ms)'] > 600]  # filter warm reqs from cold reqs
                    else:
                        data = data[data['Client Latency (ms)'] < 400]  # filter cold reqs from warm reqs

                read_latencies = data['Client Latency (ms)'].to_numpy()
                sorted_latencies = np.sort(read_latencies)
                burstsize_to_latencies[burst_size] = sorted_latencies

                plot_individual_cdf(experiment, iat_interval, sorted_latencies, burst_size)

        return burstsize_to_latencies

    title = f'{args.provider} Bursty Behavior Analysis'
    fig, axes = plt.subplots(nrows=1, ncols=2, sharey=True, figsize=(10, 5))
    fig.suptitle(title, fontsize=16)

    iat_burst_sizes_latencies = {
        '3s': plot_composing_cdf_return_latencies(axes[0], iat_interval=[3, 30], xstart=0, xend=1000),
        '600s': plot_composing_cdf_return_latencies(axes[1], iat_interval=[600, 1000], xstart=0, xend=5000)}

    plot_dual_cdf(path=args.path, latencies_dict=iat_burst_sizes_latencies, burst_size='1')
    plot_dual_cdf(path=args.path, latencies_dict=iat_burst_sizes_latencies, burst_size='500')

    plt.legend(loc='lower right')
    fig.tight_layout(rect=[0, 0, 1, 0.95])
    fig.savefig(f'{args.path}/{title}.png')
    fig.savefig(f'{args.path}/{title}.pdf')
    plt.close()

    print("Completed successfully.")
