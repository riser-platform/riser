package e2e

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ingressClient struct {
	client *http.Client
}

func NewIngressClient(ingressIP string) *ingressClient {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	return &ingressClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					// For all requests to use the ingress IP address since we don't necessarily have DNS setup for the cluster
					return dialer.DialContext(ctx, network, fmt.Sprintf("%s:443", ingressIP))
				},
				ForceAttemptHTTP2:   true,
				TLSHandshakeTimeout: 5 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *ingressClient) Get(addr string) (*http.Response, error) {
	return c.client.Get(addr)
}
