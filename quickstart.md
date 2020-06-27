# Riser Quickstart

The Riser demo is the easiest way to experiment with Riser. The demo is a single cluster environment

## Prerequisites

- You must have [git](https://git-scm.com/downloads) installed.
- You must have a GitHub account. Riser will support other git providers in the future.
- You must have [minikube](https://github.com/kubernetes/minikube) v1.4+ installed.

> :information_source: Windows Users: A Windows release is available but has not yet been tested. It's recommended that you use the [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/faq) for the Riser CLI.

### Installation

- Enable the minikube ingress addon: `minikube addons enable ingress`
- Create a minikube cluster. For the best results use the recommended settings: `minikube start --cpus=4 --memory=6144 --kubernetes-version=1.16.9`.
- Create a GitHub repo for Riser's state (e.g. https://github.com/your-name/riser-state).
- Download the [latest Riser CLI](https://github.com/riser-platform/riser/releases/) for your platform and put it in your path.
- Ensure that your minikube is started. In a new terminal window, run `minikube tunnel`. Ensure it establishes the tunnel and let it run in the backround.
- Run `riser demo install` and follow the instructions.

### Things to try

- Use `riser apps init` to create a minimal app config.
- Check out the documented [app config](examples/app.yaml) for a full list of configuration options.
- Use `riser help` to explore other help topics.
- Review the [emojivoto microservices example](examples/emojivoto) 



