package workspace

import (
	"path"

	"github.com/gravitational/gamma/internal/action"
	"github.com/gravitational/gamma/internal/node"
)

type Workspace interface {
	CollectActions() ([]action.Action, error)
}

type workspace struct {
	workingDirectory string
	outputDirectory  string
	packages         node.PackageService
}

func New(workingDirectory, outputDirectory string) Workspace {
	return &workspace{
		workingDirectory,
		outputDirectory,
		node.NewPackageService(workingDirectory),
	}
}

func (w *workspace) CollectActions() ([]action.Action, error) {
	rootPackage, err := w.readRootPackage()
	if err != nil {
		return nil, err
	}

	workspaces, err := w.packages.GetWorkspaces(rootPackage)
	if err != nil {
		return nil, err
	}

	var actions []action.Action
	for _, ws := range workspaces {
		outputDirectory := path.Join(w.outputDirectory, ws.Name)

		config := &action.Config{
			Name:             ws.Name,
			WorkingDirectory: w.workingDirectory,
			OutputDirectory:  outputDirectory,
			PackageInfo:      ws,
		}

		action, err := action.New(config)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (w *workspace) readRootPackage() (*node.PackageInfo, error) {
	p := path.Join(w.workingDirectory, "package.json")

	return w.packages.ReadPackageInfo(p)
}
