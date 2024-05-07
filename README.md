# ![design](design/logo.png)

Serverless computing has seen rapid adoption because of its instant scalability, flexible billing model, and economies of scale. In serverless,
developers structure their applications as a collection of functions invoked by various events like clicks, and cloud providers take responsibility
for cloud infrastructure management. As with other cloud services, serverless deployments require responsiveness and performance predictability manifested
through low average and tail latencies. While the average end-to-end latency has been extensively studied in prior works, existing papers lack a detailed
characterization of the effects of tail latency in real-world serverless scenarios and their root causes.

In response, we introduce STeLLAR, an open-source serverless benchmarking framework, which enables an accurate performance characterization of serverless
deployments. STeLLAR is provider-agnostic and highly configurable, allowing the analysis of both end-to-end and per-component performance with minimal
instrumentation effort. Using STeLLAR, we study three leading serverless clouds and reveal that storage accesses and bursty function invocation traffic
are key factors impacting tail latency in modern serverless systems. Finally, we identify important factors that do **not** contribute to latency variability,
such as the choice of language runtime.

## Referencing our work

If you decide to use STeLLAR for your research and experiments, we are thrilled to support you by offering
advice for potential extensions of vHive and always open for collaboration.

Please cite our [paper](docs/STeLLAR_IISWC21.pdf) that has recently been accepted to IISWC 2021:
```
@inproceedings{ustiugov:analyzing,
  author    = {Dmitrii Ustiugov and
               Theodor Amariucai and
               Boris Grot},
  title     = {Analyzing Tail Latency in Serverless Clouds with STeLLAR},
  booktitle = {Proceedings of the 2021 IEEE International Symposium on Workload Characterization (IISWC)},
  publisher = {{IEEE}},
  year      = {2021},
  doi       = {},
}
```

## Getting started with STeLLAR

STeLLAR can be readily deployed on premises or in the cloud. We provide [a quick-start guide](https://github.com/ease-lab/STeLLAR/wiki)
that describes the intial setup, as well as how to set up benchmarking experiments. More details of the STeLLAR design can be found in our IISWC'21 [paper](docs/STeLLAR_IISWC21.pdf)).


### Getting help and contributing

We would be happy to answer any questions in GitHub Issues and encourage the open-source community
to submit new Issues, assist in addressing existing issues and limitations, and contribute their code with Pull Requests.

### Continuous Benchmarking

We provide scheduled benchmarks for AWS Lambda, Azure Functions, Google Cloudrun and Cloudflare Workers daily, which can be found on our [dashboard](https://vhive-serverless.github.io/STeLLAR/). Currently, there are benchmarks for warm and cold function invocations.

## License and copyright

STeLLAR is free. We publish the code under the terms of the MIT License that allows distribution, modification, and commercial use.
This software, however, comes without any warranty or liability.

The software is maintained at the [EASE lab](https://easelab.inf.ed.ac.uk/) as part of the University of Edinburgh.


### Maintainers

* Dmitrii Ustiugov: [GitHub](https://github.com/ustiugov),
  [twitter](https://twitter.com/DmitriiUstiugov), [web page](http://homepages.inf.ed.ac.uk/s1373190/)
* [Dilina Dehigama](https://github.com/dilinade)
