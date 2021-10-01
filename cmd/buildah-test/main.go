package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/werf/werf/pkg/docker"
	"github.com/werf/werf/pkg/util"

	"github.com/werf/werf/pkg/buildah"
	"github.com/werf/werf/pkg/werf"
)

var errUsage = errors.New("./buildah-test {auto|native-rootless|docker-with-fuse} DOCKERFILE_PATH [CONTEXT_PATH]")

func do(ctx context.Context) error {
	shouldTerminate, mode, err := buildah.InitProcess(func() (buildah.Mode, error) {
		if len(os.Args) < 2 {
			return "", errUsage
		}
		return buildah.Mode(os.Args[1]), nil
	})
	if err != nil {
		return fmt.Errorf("buildah init process: %s", err)
	}
	if shouldTerminate {
		return nil
	}

	if err := werf.Init("", ""); err != nil {
		return fmt.Errorf("unable to init werf subsystem: %s", err)
	}

	if mode == buildah.ModeDockerWithFuse {
		if err := docker.Init(ctx, "", false, false, ""); err != nil {
			return err
		}
	}

	if len(os.Args) < 3 {
		return errUsage
	}

	var dockerfilePath = os.Args[2]

	var contextDir string
	if len(os.Args) > 3 {
		contextDir = os.Args[3]
	}

	b, err := buildah.NewBuildah(mode, buildah.BuildahOpts{})
	if err != nil {
		return fmt.Errorf("unable to create buildah client: %s", err)
	}

	dockerfileData, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("error reading %q: %s", dockerfilePath, err)
	}

	var contextTar io.Reader
	if contextDir != "" {
		contextTar = util.ReadDirAsTar(contextDir)
	}

	imageId, err := b.BuildFromDockerfile(ctx, dockerfileData, buildah.BuildFromDockerfileOpts{
		ContextTar: contextTar,
		CommonOpts: buildah.CommonOpts{
			LogWriter: os.Stdout,
		},
	})
	if err != nil {
		return fmt.Errorf("BuildFromDockerfile failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "INFO: built imageId is %s\n", imageId)

	return nil
}

func main() {
	if err := do(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
