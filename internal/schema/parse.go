package schema

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/gravitational/gamma/internal/cache"
)

var configCache = cache.New[*Config]()

func GetConfig(root, filename string) (*Config, error) {
	var config CustomConfig

	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filename, err)
	}

	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", filename, err)
	}

	config.Path = filename

	return parseCustomConfig(root, filename, config)
}

func parseCustomConfig(root, filename string, customConfig CustomConfig) (*Config, error) {
	config := &Config{
		Path:        customConfig.Path,
		Name:        customConfig.Name,
		Author:      customConfig.Author,
		Description: customConfig.Description,
		Inputs:      customConfig.Inputs,
		Outputs:     customConfig.Outputs,
		Runs:        customConfig.Runs,
		Branding:    customConfig.Branding,
	}

	if customConfig.Extend != nil {
		for _, extension := range *customConfig.Extend {
			file := extension.From
			if strings.HasPrefix(file, "@/") {
				file = strings.TrimPrefix(file, "@/")
				file = path.Join(root, file)
			}
			if !path.IsAbs(file) {
				file = path.Join(filename, file)
			}

			var extensionConfig *Config
			var ok bool

			extensionConfig, ok = configCache.Get(file)
			if !ok {
				def, err := GetConfig(root, file)
				if err != nil {
					return nil, err
				}

				configCache.Set(file, def)

				extensionConfig = def
			}

			if err := mergeConfigs(config, extensionConfig, extension.Include); err != nil {
				return nil, err
			}
		}
	}

	return config, nil
}

func mergeConfigs(base, extension *Config, includes *[]ExtensionInclude) error {
	if includes != nil {
		for _, include := range *includes {
			field := include.Field

			switch field {
			case "inputs":
				if err := mergeInputs(base, extension, &include); err != nil {
					return err
				}
			case "outputs":
				if err := mergeOutputs(base, extension, &include); err != nil {
					return err
				}
			case "branding":
				if err := mergeBranding(base, extension, &include); err != nil {
					return err
				}
			case "author":
				if err := mergeAuthor(base, extension); err != nil {
					return err
				}
			case "runs":
				if err := mergeRuns(base, extension); err != nil {
					return err
				}
			}
		}

		return nil
	}

	_ = mergeInputs(base, extension, nil)
	_ = mergeOutputs(base, extension, nil)
	_ = mergeBranding(base, extension, nil)
	_ = mergeAuthor(base, extension)
	_ = mergeRuns(base, extension)

	return nil
}

func mergeInputs(base, extension *Config, includes *ExtensionInclude) error {
	if extension.Inputs == nil {
		return fmt.Errorf("no inputs exist in %s", extension.Path)
	}

	newInputs := make(InputMap)

	if base.Inputs != nil {
		for key, value := range *base.Inputs {
			newInputs[key] = value
		}
	}

	if includes != nil && includes.Include != nil {
		for _, field := range *includes.Include {

			inputs := *extension.Inputs

			input, ok := inputs[field]
			if !ok {
				return fmt.Errorf("input %s does not exist in %s", field, extension.Path)
			}

			if _, ok := newInputs[field]; ok {
				return fmt.Errorf("conflicting input, %s exists in both %s and %s", field, base.Path, extension.Path)
			}

			newInputs[field] = input
		}
	} else {
	outer:
		for key, input := range *extension.Inputs {
			if includes != nil && includes.Exclude != nil {
				for _, exclude := range *includes.Exclude {
					if exclude == key {
						continue outer
					}
				}
			}

			newInputs[key] = input
		}
	}

	base.Inputs = &newInputs

	return nil
}

func mergeOutputs(base, extension *Config, includes *ExtensionInclude) error {
	if extension.Outputs == nil {
		return fmt.Errorf("no outputs exist in %s", extension.Path)
	}

	newOutputs := make(OutputMap)

	if base.Outputs != nil {
		for key, value := range *base.Outputs {
			newOutputs[key] = value
		}
	}

	if includes != nil && includes.Include != nil {
		for _, field := range *includes.Include {

			outputs := *extension.Outputs

			output, ok := outputs[field]
			if !ok {
				return fmt.Errorf("output %s does not exist in %s", field, extension.Path)
			}

			if _, ok := newOutputs[field]; ok {
				return fmt.Errorf("conflicting output, %s exists in both %s and %s", field, base.Path, extension.Path)
			}

			newOutputs[field] = output
		}
	} else {
	outer:
		for key, output := range *extension.Outputs {
			if includes != nil && includes.Exclude != nil {
				for _, exclude := range *includes.Exclude {
					if exclude == key {
						continue outer
					}
				}
			}

			newOutputs[key] = output
		}
	}

	base.Outputs = &newOutputs

	return nil
}

func mergeBranding(base, extension *Config, includes *ExtensionInclude) error {
	if extension.Branding == nil {
		return fmt.Errorf("no branding exists in %s", extension.Path)
	}

	newBranding := &Branding{}
	if base.Branding != nil {
		newBranding = base.Branding
	}

	if includes != nil && includes.Include != nil {
		for _, field := range *includes.Include {

			switch field {
			case "color":
				newBranding.Color = extension.Branding.Color
			case "icon":
				newBranding.Icon = extension.Branding.Icon
			}
		}
	} else {
		var excludeColor bool
		var excludeIcon bool

		if includes != nil && includes.Exclude != nil {
			for _, exclude := range *includes.Exclude {
				switch exclude {
				case "color":
					excludeColor = true
				case "icon":
					excludeIcon = true
				}
			}
		}

		if !excludeColor {
			newBranding.Color = extension.Branding.Color
		}

		if !excludeIcon {
			newBranding.Icon = extension.Branding.Icon
		}
	}

	base.Branding = newBranding

	return nil
}

func mergeRuns(base, extension *Config) error {
	if extension.Runs.JavascriptRun == nil &&
		extension.Runs.DockerRun == nil &&
		extension.Runs.CompositeRun == nil {
		return fmt.Errorf("runs is empty in %s", extension.Path)
	}

	if extension.Runs.JavascriptRun != nil {
		base.Runs = Runs{
			JavascriptRun: extension.Runs.JavascriptRun,
		}
	}

	if extension.Runs.DockerRun != nil {
		base.Runs = Runs{
			DockerRun: extension.Runs.DockerRun,
		}
	}

	if extension.Runs.CompositeRun != nil {
		base.Runs = Runs{
			CompositeRun: extension.Runs.CompositeRun,
		}
	}

	return nil
}

func mergeAuthor(base, extension *Config) error {
	if extension.Author == nil {
		return fmt.Errorf("author is empty in %s", extension.Path)
	}

	base.Author = extension.Author

	return nil
}
