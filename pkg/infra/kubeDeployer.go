package infra

type KubeDeployer interface {
	Deploy() error
}

// NoopDeployer does nothing. It's intended to be used when existing k8s infra is in place.
type NoopDeployer struct{}

func (*NoopDeployer) Deploy() error {
	return nil
}
