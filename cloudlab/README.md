This folder contains the necessary scripts to set up/stabilize a clean CloudLab experiment for fccd-orchestrator. Information about choosing the right profile and nodes on CloudLab can be found in the [CloudLab Guidelines Wiki](https://github.com/ustiugov/fccd-orchestrator/wiki/CloudLab-Guidelines). 

## Prerequisites
1. Enter your public key in GitHub.
2. Configure your email and user for Git.

## Stabilization
1. `./stabilize.sh`
2. Start a new login shell,
3. `./setup_containerd.sh`

You can now start firecracker-containerd, clone the fccd-orchestrator repo, and start developing or testing.
