package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type Workspaces struct {
	Value []string
}

func (w *Workspaces) UnmarshalJSON(data []byte) error {
	var obj struct {
		Packages []string `json:"packages"`
	}

	if err := json.Unmarshal(data, &obj); err == nil {
		w.Value = obj.Packages

		return nil
	}
	var array []string
	if err := json.Unmarshal(data, &array); err == nil {
		w.Value = array

		return nil
	}

	return errors.New("could not parse workspaces")
}

type packageService struct {
	RootPath string
}

type PackageService interface {
	ReadPackageInfo(filename string) (*PackageInfo, error)
	GetWorkspaces(p *PackageInfo) ([]*PackageInfo, error)
}

func NewPackageService(rootPath string) PackageService {
	return &packageService{
		rootPath,
	}
}

func (s *packageService) ReadPackageInfo(filename string) (*PackageInfo, error) {
	p := PackageInfo{
		Path:     path.Dir(filename),
		RootPath: s.RootPath,
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %v [%s]", err, filename)
	}

	if err := json.Unmarshal(contents, &p); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v [%s]", err, filename)
	}

	return &p, nil
}

type PackageInfo struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Repository string     `json:"repository"`
	Workspaces Workspaces `json:"workspaces"`

	Path     string
	RootPath string
}

func (s *packageService) GetWorkspaces(p *PackageInfo) ([]*PackageInfo, error) {
	if len(p.Workspaces.Value) == 0 {
		return nil, errors.New("no workspaces specified")
	}

	var workspaces []*PackageInfo

	dir := os.DirFS(p.Path)
	for _, workspace := range p.Workspaces.Value {
		matches, err := fs.Glob(dir, workspace)

		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			filename := path.Join(p.Path, match, "package.json")

			w, err := s.ReadPackageInfo(filename)
			if err != nil {
				return nil, err
			}

			workspaces = append(workspaces, w)
		}
	}

	return workspaces, nil
}
