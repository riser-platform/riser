# Riser

<p align="left">
    <a href="https://github.com/riser-platform/riser"><img alt="GitHub Actions status" src="https://github.com/riser-platform/riser/workflows/Build/badge.svg"></a>
</p>

Riser is an opinionated app platform built on Kubernetes. It provides a radically simplified application deployment and management experience. without vendor lock-in.

[![asciicast](https://asciinema.org/a/277448.svg)](https://asciinema.org/a/277448?autoplay=1&cols=160&rows=40)

> :warning: This is an experimental project with the goal of improving how we deploy and manage common application workloads. You're invited to look around and provide feedback. Until there is a stable release it is not advised to use Riser in production as breaking changes may occur frequently and without warning.

**[Check out the quickstart!](quickstart.md)**

## Key Features

- Radically simplified deployment and management of [12 factor apps](https://12factor.net/)
- 100% Open Source. PaaS experience without vendor or cloud lock-in
- Single view of apps across multiple stages (cluster)
- Simplified secrets management
- Loosely coupled. Thanks to a purely GitOps approach, Riser can go down or even be completely removed without affecting your apps
- Developers only need access to Riser. Kubernetes access is optional for advanced debugging or operational tasks

**[Check out the quickstart!](quickstart.md)**

### More on GitOps

Riser interacts with Kubernetes using a strictly [GitOps](https://thenewstack.io/what-is-gitops-and-why-it-might-be-the-next-big-thing-for-devops/) approach. A git repository (typically referred to as a "state repo") contains all information required to stand up an app. The riser server can be unreachable or even destroyed with no impact to your apps. It also designed so that all of the resources in your state repo can be "`kubectl apply -f`'d" to a Kubernetes cluster without any Riser infrastructure running or installed.

## Known Issues and Limitations

- Docker images must be hosted on a public registry. Private registries as well as the ability for administrators to restrict the use of public registries will be supported in the future.
- While not strictly enforced, GitHub is the only supported git provider at this time.
- The documentation is very sparse. As features mature more documentation will be added.

## Development

We are not currently accepting PRs. As the project matures this section will contain more details.

### Assets

If you change anything in the `assets` folder, you must run `make generate` to statically bundle them inside the riser binary.

### Related projects

- [Riser Server](https://github.com/riser-platform/riser-server)
- [Riser Controller](https://github.com/riser-platform/riser-controller)
