package tsc

import (
	"context"
	"os"
	"sync"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/confile"
	"github.com/tendermint/starport/starport/pkg/nodetime"
)

// Config represents tsconfig.json.
type Config struct {
	Include         []string          `json:"include"`
	CompilerOptions []CompilerOptions `json:"compilerOptions"`
}

// CompilerOptions section of tsconfig.json.
type CompilerOptions struct {
	Declaration string              `json:"declaration"`
	Paths       map[string][]string `json:"paths"`
}

var placeOnce sync.Once

// Generate transpiles TS into JS by given TS config.
func Generate(ctx context.Context, config Config) error {
	var err error

	placeOnce.Do(func() { err = nodetime.PlaceBinary() })

	if err != nil {
		return err
	}

	// save the config into a temp file in the fs.
	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	f.Close()
	defer os.Remove(f.Name())

	if err := confile.
		New(confile.DefaultJSONEncodingCreator, f.Name()).
		Save(config); err != nil {
		return err
	}

	// command constructs the tsc command.
	command := []string{
		nodetime.BinaryPath,
		nodetime.CommandTSC,
		"-b",
		f.Name(),
	}

	// execute the command.
	return cmdrunner.Exec(ctx, command[0], command[1:]...)
}
