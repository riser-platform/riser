package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

const (
	// DefaultKindNodeImage should roughly match the latest stable kubernetes version provided by GKE/AKS/EKS
	DefaultKindNodeImage = "kindest/node:v1.16.9"
	DefaultKindName      = "riser-e2e"
)

func main() {
	var kindNodeImage string
	flag.StringVar(&kindNodeImage, "name", DefaultKindNodeImage, "help message for flagname")

	args := []string{"create", "cluster", fmt.Sprintf("--image=%s", kindNodeImage), fmt.Sprintf("--name=%s", DefaultKindName)}

	cmd := exec.Command("kind", args...)
	output, _ := cmd.CombinedOutput()

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing kind: %v", string(output))
	}
}
