package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Path        string     `yaml:"-"`
	Name        string     `yaml:"name"`
	Author      *string    `yaml:"author,omitempty"`
	Description string     `yaml:"description"`
	Inputs      *InputMap  `yaml:"inputs,omitempty"`
	Outputs     *OutputMap `yaml:"output,omitempty"`
	Runs        Runs       `yaml:"runs"`
	Branding    *Branding  `yaml:"branding,omitempty"`
}

type CustomConfig struct {
	Path        string       `yaml:"-"`
	Name        string       `yaml:"name"`
	Author      *string      `yaml:"author,omitempty"`
	Description string       `yaml:"description"`
	Inputs      *InputMap    `yaml:"inputs,omitempty"`
	Outputs     *OutputMap   `yaml:"output,omitempty"`
	Runs        Runs         `yaml:"runs"`
	Branding    *Branding    `yaml:"branding,omitempty"`
	Extend      *[]Extension `yaml:"extend,omitempty"`
}

type ExtensionInclude struct {
	Field   string    `yaml:"field"`
	Include *[]string `yaml:"include"`
	Exclude *[]string `yaml:"exclude"`
}

type Extension struct {
	From    string              `yaml:"from"`
	Include *[]ExtensionInclude `yaml:"include"`
}

type Branding struct {
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}

type Input struct {
	Description        string  `yaml:"description"`
	Required           *bool   `yaml:"required,omitempty"`
	Default            *string `yaml:"default,omitempty"`
	DeprecationMessage *string `yaml:"deprecationMessage,omitempty"`
}

type InputMap = map[string]Input

type Output struct {
	Description string `yaml:"description"`
	Value       string `yaml:"value"`
}

type OutputMap = map[string]Output

type EnvMap = map[string]string
type WithMap = map[string]string

type JavascriptRun struct {
	Using  string  `yaml:"using"`
	Main   string  `yaml:"main"`
	Pre    *string `yaml:"pre,omitempty"`
	PreIf  *string `yaml:"pre-if,omitempty"`
	Post   *string `yaml:"post,omitempty"`
	PostIf *string `yaml:"post-if,omitempty"`
}

type RunStep struct {
	Run              *string  `yaml:"run,omitempty"`
	Shell            *string  `yaml:"shell,omitempty"`
	If               *string  `yaml:"if,omitempty"`
	Name             *string  `yaml:"name,omitempty"`
	ID               *string  `yaml:"id,omitempty"`
	Env              *EnvMap  `yaml:"env,omitempty"`
	WorkingDirectory *string  `yaml:"working-directory,omitempty"`
	Uses             *string  `yaml:"uses,omitempty"`
	With             *WithMap `yaml:"with,omitempty"`
}

type CompositeRun struct {
	Using string    `yaml:"using"`
	Steps []RunStep `yaml:"steps"`
}

type DockerRun struct {
	PreEntrypoint  *string   `yaml:"pre-entrypoint,omitempty"`
	Image          string    `yaml:"image"`
	Env            *EnvMap   `yaml:"env,omitempty"`
	Entrypoint     *string   `yaml:"entrypoint,omitempty"`
	PostEntrypoint *string   `yaml:"post-entrypoint,omitempty"`
	Args           *[]string `yaml:"args,omitempty"`
}

type Runs struct {
	Using string `yaml:"using"`

	*CompositeRun
	*JavascriptRun
	*DockerRun
}

func (r Runs) MarshalYAML() (interface{}, error) {
	if r.JavascriptRun != nil {
		return r.JavascriptRun, nil
	}

	if r.DockerRun != nil {
		return r.DockerRun, nil
	}

	if r.CompositeRun != nil {
		return r.CompositeRun, nil
	}

	return nil, nil
}

func (r *Runs) UnmarshalYAML(value *yaml.Node) error {
	var obj struct {
		Using string `yaml:"using"`
	}
	if err := value.Decode(&obj); err != nil {
		return err
	}

	r.Using = obj.Using

	switch obj.Using {
	case "composite":
		var compositeRun CompositeRun

		if err := value.Decode(&compositeRun); err != nil {
			return err
		}

		r.CompositeRun = &compositeRun

		return nil

	case "docker":
		var dockerRun DockerRun

		if err := value.Decode(&dockerRun); err != nil {
			return err
		}

		r.DockerRun = &dockerRun

		return nil

	case "node12", "node16":
		var javascriptRun JavascriptRun

		if err := value.Decode(&javascriptRun); err != nil {
			return err
		}

		r.JavascriptRun = &javascriptRun

		return nil
	}

	return fmt.Errorf("unsupported runs.using value: %v, expected composite, docker, node12 or node16", obj.Using)
}
