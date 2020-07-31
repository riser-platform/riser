# Riser

<p align="left">
    <a href="https://github.com/riser-platform/riser"><img alt="GitHub Actions status" src="https://github.com/riser-platform/riser/workflows/Build/badge.svg"></a>
</p>

Riser is an opinionated app platform built on [Kubernetes](https://kubernetes.io/) and [Knative](https://knative.dev). It provides a radically simplified application deployment and management experience without vendor lock-in.

[![asciicast](https://asciinema.org/a/350860.svg)](https://asciinema.org/a/350860?autoplay=1&cols=160&rows=40)

> :warning: This is an experimental project with the goal of improving how we deploy and manage common application workloads. You're invited to look around and provide feedback. It is not yet advised to use Riser in production. Breaking changes may occur frequently and without warning.

## Key Features

- Radically simplified deployment and management of [12 factor apps](https://12factor.net/)
- PaaS experience without vendor or cloud lock-in
- Single view of apps across multiple environments (e.g. dev/test/prod)
- Simplified secrets management
- GitOps: All state changes happen through git
- App developers only need access to Riser. Kubernetes cluster access is optional for advanced debugging or operational tasks

**[Check out the quickstart!](https://docs.riser.dev/docs/quickstart/)**

## Development

> Note: We are not currently accepting PRs. As the project matures this section will contain more details.

### Assets

If you change anything in the `assets` folder, you must run `make generate` to statically bundle them inside the riser binary.

### E2E Tests using Kind

#### Prequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- `envsubst` (`gettext` package for homebrew, apt, etc.)
- A git repo with write access

#### Building the E2E image

- Check out this repo to the desired tag that you wish to test
- Run `make docker-e2e`

#### Running

Example using a github deploy key with write access:
```
go run pkg/e2e/kind/main.go --git-url git@github.com:me/riser-state --git-ssh-key-path=/Users/me/.ssh/id_rsa
```

Run `go run pkg/e2e/kind/main.go` for additional options


### Supporting projects

- [Riser Server](https://github.com/riser-platform/riser-server)
- [Riser Controller](https://github.com/riser-platform/riser-controller)
