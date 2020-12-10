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


def plot_cpu_slowdown():
    def add_memory_result_to_plots(memory_allocated_mb, color):
        latencies, service_times = get_experiment_results(memory_allocated_mb)

        # plot latencies from experiments
        plt.plot(service_times, latencies, 'o',
                 color=color)  # , label='recorded latencies (' + memory_allocated_mb + 'MB)'
        time = np.linspace(0, 1250, 50)

        # aws_latency = amp x service_time + infra_const
        latencies -= np.array(fixed_infra_constant)
        res = np.linalg.lstsq(service_times, latencies, rcond=None)[0]
        amplification = res[0][0]
        cpu_utilization = 1 / amplification
        memory_to_util[memory_allocated_mb] = cpu_utilization

        print("Amplification for {0}MB: {1:0.2f}".format(memory_allocated_mb, amplification))
        print("CPU utilization for {0}MB: {1:0.2f}".format(memory_allocated_mb, cpu_utilization))

        pred_latency = res[0] * time + fixed_infra_constant
        plt.plot(time, pred_latency, color=color, label='predicted latency (' + memory_allocated_mb + 'MB)')

    def get_experiment_results(memory_allocated_mb):
        experiment_dirs = []
        for dir_path, dir_names, filenames in os.walk(path):
            # filter by memory allocated and no subdirectories
            if not dir_names and dir_path.split('/')[2].split('MB')[0] == memory_allocated_mb:
                experiment_dirs.append(dir_path)

        latencies = np.empty((0, 1), float)
        service_times = np.empty((0, 1), float)
        for experiment in experiment_dirs:
            experiment_name = experiment.split('/')[2]
            service_time_sec = int(experiment_name.split('-')[2].split('ms')[0])
            with open(experiment + "/latencies.csv") as file:
                data = pd.read_csv(file)
                read_latencies = data['Client Latency (ms)'].to_numpy()
                read_latencies_no = len(read_latencies)
                latencies = np.vstack((latencies, read_latencies.reshape((read_latencies_no, 1))))
                service_times = np.vstack((service_times, np.ones((read_latencies_no, 1)) * service_time_sec))

        return latencies, service_times

    title = provider + ' CPU Slowdown'
    fig = plt.figure(figsize=(12, 5))
    fig.suptitle(title)
    plt.ylabel('Latency (ms)')
    plt.xlabel('Service Time (ms)')

    memory_to_util = dict({})

    path = 'providers/{0}/'.format(provider)
    add_memory_result_to_plots('128', 'tab:blue')
    add_memory_result_to_plots('480', 'tab:orange')
    add_memory_result_to_plots('832', 'tab:green')
    add_memory_result_to_plots('1184', 'tab:pink')
    add_memory_result_to_plots('1536', 'tab:purple')
    add_memory_result_to_plots('2304', 'tab:red')
    add_memory_result_to_plots('5120', 'tab:grey')
    add_memory_result_to_plots('10240', 'tab:cyan')

    plt.legend(loc='upper left')
    fig.savefig(path + title + '.png')
    plt.close()

    title2 = provider + ' CPU Utilization Rates'
    fig = plt.figure(figsize=(5, 5))
    fig.suptitle(title2)
    axes = plt.gca()
    axes.set_ylim([0, 1])
    # axes.set_xlim([0, 1536])
    plt.xlabel('Function Memory (MB)')
    plt.ylabel('Fraction')
    for memory_mb in memory_to_util:
        plt.plot(memory_mb, memory_to_util[memory_mb], 'o', color='tab:blue')
    fig.savefig('providers/' + provider + '/' + title2 + '.png')
    plt.close()


fixed_infra_constant = 50
print("Infrastructure constant: {0:0.2f}ms".format(fixed_infra_constant))

# First run experiment `cpu-slowdown`
# Before running this plotting utility, place your results under, e.g., for
# AWS, providers/AWS.
provider = 'AWS'
plot_cpu_slowdown()
print("Completed successfully.")
