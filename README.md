# Riser

<p align="left">
    <a href="https://github.com/riser-platform/riser"><img alt="GitHub Actions status" src="https://github.com/riser-platform/riser/workflows/Main/badge.svg"></a>
</p>

Riser is an opinionated app platform built on Kubernetes. It provides radically simplified application deployment and management without vendor lock-in.

> :warning: This is an experimental project that is in the proof of concept phase. Breaking changes may occur frequently. You're more than welcome to look around and provide feedback, but this project is not recommended for mission critical workloads.

**[Check out the quickstart!](quickstart.md)**

## Key Features

- Radically simplified deployment and management of [12 factor apps](https://12factor.net/)
- 100% Open Source. PaaS experience without vendor or cloud lock-in
- Single view of apps across multiple stages (cluster)
- Simplified secrets management
- Loosely coupled. Thanks to a purely GitOps approach, Riser can go down or even be completely removed your workloads continue running
- Developers only need access to Riser. Kubernetes access is optional for advanced debugging or operational tasks

**[Check out the quickstart!](quickstart.md)**

### GitOps

Riser interacts with Kubernetes using a strictly [GitOps](https://thenewstack.io/what-is-gitops-and-why-it-might-be-the-next-big-thing-for-devops/) approach. A git repository (typically referred to as a "state repo") contains all information required to stand up an app. The riser server can be unreachable or even destroyed with no impact to an app. Additionally, Riser does not require any of its own custom types (CRDs) for app. This means that all of the resources in your state repo can be "`kubectl apply -f`'d" to a Kubernetes cluster without Riser and operate as desired. Of course, any supporting infrastructure such as Istio must still be configured.

## Known Issues and Limitations

- Docker images must be hosted on a public registry. Private registries as well as the ability for administrators to restrict the use of public registries will be supported in the future.
- While not strictly enforced, GitHub is the only supported git provider at this time.
- The documentation is very sparse. As features mature more documentation will be added.

## Development

We are not currently accepting PRs. As the project matures this section will contain more details.

### Assets

If you change anything in the `assets` folder, you must run `make generate` to statically bundle them inside the riser binary.

### Related projects

TODO: link to server and controller.