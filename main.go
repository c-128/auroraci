package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/docker/docker/client"
)

func main() {
	auroraDirectory, found := os.LookupEnv("AURORACI_DIRECTORY")
	if !found {
		panic("missing AURORA_DIRECTORY environment variable")
	}

	pipelinesDirectory := path.Join(auroraDirectory, "pipelines")
	artifactsDirectory := path.Join(auroraDirectory, "artifacts")

	err := errors.Join(
		os.MkdirAll(pipelinesDirectory, os.ModePerm),
		os.MkdirAll(artifactsDirectory, os.ModePerm),
	)
	if err != nil {
		panic(err)
	}

	// DOCKER_HOST environment variable
	docker, err := client.NewClientWithOpts(client.WithHostFromEnv())
	if err != nil {
		panic(err)
	}

	fmt.Println(docker)
}
