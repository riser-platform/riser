package infra

type KubeDeployer interface {
	Deploy() error
	Destroy() error
}

// NoopDeployer does nothing. It's intended to be used when existing k8s infra is in place.
type NoopDeployer struct{}

func (*NoopDeployer) Deploy() error {
	return nil
}

func (*NoopDeployer) Destroy() error {
	return nil
}
