package instruction

import (
	"context"

	"github.com/werf/werf/pkg/build/stage"
	"github.com/werf/werf/pkg/config"
	"github.com/werf/werf/pkg/container_backend"
	backend_instruction "github.com/werf/werf/pkg/container_backend/instruction"
	"github.com/werf/werf/pkg/dockerfile"
	dockerfile_instruction "github.com/werf/werf/pkg/dockerfile/instruction"
	"github.com/werf/werf/pkg/util"
)

type StopSignal struct {
	*Base[*dockerfile_instruction.StopSignal]
}

func NewStopSignal(name stage.StageName, i *dockerfile.DockerfileStageInstruction[*dockerfile_instruction.StopSignal], dependencies []*config.Dependency, hasPrevStage bool, opts *stage.BaseStageOptions) *StopSignal {
	return &StopSignal{Base: NewBase(name, i, backend_instruction.NewStopSignal(*i.Data), dependencies, hasPrevStage, opts)}
}

func (stage *StopSignal) GetDependencies(ctx context.Context, c stage.Conveyor, cb container_backend.ContainerBackend, prevImage, prevBuiltImage *stage.StageImage, buildContextArchive container_backend.BuildContextArchiver) (string, error) {
	var args []string
	args = append(args, stage.instruction.Data.Name())
	args = append(args, stage.instruction.Data.Signal)
	return util.Sha256Hash(args...), nil
}
