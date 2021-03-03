package cosmosgen

import (
	"context"

	"github.com/tendermint/starport/starport/pkg/cmdrunner"
	"github.com/tendermint/starport/starport/pkg/cmdrunner/step"
	"github.com/tendermint/starport/starport/pkg/cosmosanalysis/module"
	"github.com/tendermint/starport/starport/pkg/gomodule"
	gomodmodule "golang.org/x/mod/module"
)

type generateOptions struct {
	includeDirs         []string
	gomodPath           string
	jsOut               func(module.Module) string
	jsIncludeThirdParty bool
}

// TODO add WithInstall.

// Option configures code generation.
type Option func(*generateOptions)

// WithJSGeneration adds JS code generation. out hook is called for each module to
// retrieve the path that should be used to place generated js code inside for a given module.
// if includeThirdPartyModules set to true, code generation will be made for the 3rd party modules
// used by the app -including the SDK- as well.
func WithJSGeneration(includeThirdPartyModules bool, out func(module.Module) (path string)) Option {
	return func(o *generateOptions) {
		o.jsOut = out
		o.jsIncludeThirdParty = includeThirdPartyModules
	}
}

// WithGoGeneration adds Go code generation.
func WithGoGeneration(gomodPath string) Option {
	return func(o *generateOptions) {
		o.gomodPath = gomodPath
	}
}

// IncludeDirs configures the third party proto dirs that used by app's proto.
// relative to the projectPath.
func IncludeDirs(dirs []string) Option {
	return func(o *generateOptions) {
		o.includeDirs = dirs
	}
}

// generator generates code for sdk and sdk apps.
type generator struct {
	ctx      context.Context
	appPath  string
	protoDir string
	o        *generateOptions
	deps     []gomodmodule.Version
}

// Generate generates code from protoDir of an SDK app residing at appPath with given options.
// protoDir must be relative to the projectPath.
func Generate(ctx context.Context, appPath, protoDir string, options ...Option) error {
	g := &generator{
		ctx:      ctx,
		appPath:  appPath,
		protoDir: protoDir,
		o:        &generateOptions{},
	}

	for _, apply := range options {
		apply(g.o)
	}

	if err := g.setup(); err != nil {
		return err
	}

	if g.o.gomodPath != "" {
		if err := g.generateGo(); err != nil {
			return err
		}
	}

	// js generation requires Go types to be existent in the source code.
	// so it needs to run after Go code gen.
	if g.o.jsOut != nil {
		if err := g.generateJS(); err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) setup() (err error) {
	// Cosmos SDK hosts proto files of own x/ modules and some third party ones needed by itself and
	// blockchain apps. Generate should be aware of these and make them available to the blockchain
	// app that wants to generate code for its own proto.
	//
	// blockchain apps may use different versions of the SDK. following code first makes sure that
	// app's dependencies are download by 'go mod' and cached under the local filesystem.
	// and then, it determines which version of the SDK is used by the app and what is the absolute path
	// of its source code.
	if err := cmdrunner.
		New(cmdrunner.DefaultWorkdir(g.appPath)).
		Run(g.ctx, step.New(step.Exec("go", "mod", "download"))); err != nil {
		return err
	}

	// parse the go.mod of the app and extract dependencies.
	modfile, err := gomodule.ParseAt(g.appPath)
	if err != nil {
		return err
	}

	g.deps, err = gomodule.ResolveDependencies(modfile)

	return
}
