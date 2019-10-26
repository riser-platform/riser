# Riser Quickstart

The Riser demo is the easiest way to experiment with Riser. The demo is a single cluster environment

## Prerequisites

- You must have [git](https://git-scm.com/downloads) installed.
- You must have a GitHub account. Riser will support other git providers in the future.
- You must have [minikube](https://github.com/kubernetes/minikube) v1.4+ installed.

> :information_source: Windows Users: A Windows release is available but has not yet been tested. It's recommended that you use the [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/faq) for the Riser CLI.

### Installation

- Enable the minikube ingress addon: `minikube addons enable ingress`
- Create a minikube cluster. For the best results use the recommended settings: `minikube start --cpus=4 --memory=4096 --kubernetes-version=1.14.6`.
- Once created start `minikube tunnel` in a separate terminal.
- Create a GitHub repo for Riser's state (e.g. https://github.com/your-name/riser-state).
- Download the [latest Riser CLI](https://github.com/riser-platform/riser/releases/) for your platform and put it in your path.
- Run `riser demo install` and follow the instructions.

### Things to try

- Check out the documented [app config](examples/app.yaml) for a full list of configuration options.
- Use `riser help` to explore other help topics.
- To see how app dependencies work, review the [emojivoto example](examples/emojivoto/README.md)



