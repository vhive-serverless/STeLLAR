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


def plot_cpu_stats(args):
    fixed_infra_constant = 50
    print("Infrastructure constant: {0:0.2f}ms".format(fixed_infra_constant))

    memory_to_slowdown = dict({})

    def add_memory_result_to_plots(memory_allocated_mb):
        latencies, service_times = get_experiment_results(memory_allocated_mb)

        p = plt.plot(service_times, latencies, 'o')  # , label='recorded latencies (' + memory_allocated_mb + 'MB)'
        time = np.linspace(0, 1250, 50)

        # aws_latency = amp x service_time + infra_const
        latencies -= np.array(fixed_infra_constant)
        res = np.linalg.lstsq(service_times, latencies, rcond=None)[0]
        slowdown = res[0][0]

        memory_to_slowdown[memory_allocated_mb] = slowdown
        print("Slowdown for {0}MB: {1:0.2f}".format(memory_allocated_mb, slowdown))

        pred_latency = res[0] * time + fixed_infra_constant
        plt.plot(time, pred_latency, color=p[-1].get_color(), label=f'Predicted ({memory_allocated_mb}MB)')

    def get_experiment_results(memory_allocated_mb):
        experiment_dirs = []
        for dir_path, dir_names, filenames in os.walk(args.path):
            # filter by memory allocated and no subdirectories
            if not dir_names and dir_path.split('/')[-1].split('memory')[1].split('MB')[0] == memory_allocated_mb:
                experiment_dirs.append(dir_path)

        latencies = np.empty((0, 1), float)
        service_times = np.empty((0, 1), float)
        for experiment in experiment_dirs:
            experiment_name = experiment.split('/')[-1]
            service_time_sec = int(experiment_name.split('-st')[1].split('ms')[0])
            with open(experiment + "/latencies.csv") as file:
                data = pd.read_csv(file)
                read_latencies = data['Client Latency (ms)'].to_numpy()
                read_latencies_no = len(read_latencies)
                latencies = np.vstack((latencies, read_latencies.reshape((read_latencies_no, 1))))
                service_times = np.vstack((service_times, np.ones((read_latencies_no, 1)) * service_time_sec))

        return latencies, service_times

    def plot_cpu_slowdown():
        title = f'{args.provider} CPU Slowdown'
        fig = plt.figure(figsize=(10, 5))
        fig.suptitle(title)
        plt.ylabel('Latency (ms)')
        plt.xlabel('Service Time (ms)')
        plt.grid(True)

        for memory in ['128', '480', '832', '1184', '1536', '2304', '5120', '10240']:
            add_memory_result_to_plots(memory)

        plt.legend(loc='upper left')
        fig.savefig(f'{args.path}/{title}.png')
        fig.savefig(f'{args.path}/{title}.pdf')
        plt.close()

    def plot_cpu_utilization():
        title = f'{args.provider} CPU Utilization Rates'
        fig = plt.figure(figsize=(5, 5))
        fig.suptitle(title)
        plt.xlabel('Function Memory (MB)')
        plt.ylabel('Fraction')
        plt.grid(True)

        axes = plt.gca()
        axes.set_ylim([0, 1])
        axes.set_xscale('linear')

        for memory_mb in memory_to_slowdown:
            cpu_utilization = 1 / memory_to_slowdown[memory_mb]
            plt.plot(int(memory_mb), cpu_utilization, 'o--', color='tab:blue')
            plt.annotate(f'{cpu_utilization:0.2f}', (int(memory_mb), cpu_utilization+0.01))
            print("CPU utilization for {0}MB: {1:0.2f}".format(memory_mb, cpu_utilization))

        fig.savefig(f'{args.path}/{title}.png')
        fig.savefig(f'{args.path}/{title}.pdf')
        plt.close()

    plot_cpu_slowdown()
    plot_cpu_utilization()

    print("Completed successfully.")
