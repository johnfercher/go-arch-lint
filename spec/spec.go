package spec

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

const supportedVersion = 1

type (
	yamlVendorName       = string
	yamlComponentName    = string
	yamlExcludeLocalPath = string

	YamlSpec struct {
		Version            int                                 `yaml:"version"`
		Allow              YamlAllow                           `yaml:"allow"`
		Vendors            map[yamlVendorName]YamlVendor       `yaml:"vendors"`
		Exclude            []yamlExcludeLocalPath              `yaml:"exclude"`
		ExcludeFilesRegExp []string                            `yaml:"excludeFiles"`
		Components         map[yamlComponentName]YamlComponent `yaml:"components"`
		Dependencies       map[yamlComponentName]YamlRules     `yaml:"deps"`
		CommonComponents   []yamlComponentName                 `yaml:"commonComponents"`
		CommonVendors      []yamlVendorName                    `yaml:"commonVendors"`
	}

	YamlAllow struct {
		DepOnAnyVendor bool `yaml:"depOnAnyVendor"`
	}

	YamlVendor struct {
		ImportPath string `yaml:"in"`
	}

	YamlComponent struct {
		LocalPath string `yaml:"in"`
	}

	YamlRules struct {
		MayDependOn    []yamlComponentName `yaml:"mayDependOn"`
		CanUse         []yamlVendorName    `yaml:"canUse"`
		AnyProjectDeps bool                `yaml:"anyProjectDeps"`
		AnyVendorDeps  bool                `yaml:"anyVendorDeps"`
	}

	YamlParseError struct {
		Err      error
		Warnings []Warning
	}
)

func newSpec(archFile string, rootDirectory string) (YamlSpec, YamlParseError) {
	spec := YamlSpec{}

	data, err := ioutil.ReadFile(archFile)
	if err != nil {
		return spec, YamlParseError{
			Err:      fmt.Errorf("can`t open '%s': %v", archFile, err),
			Warnings: []Warning{},
		}
	}

	reader := bytes.NewBuffer(data)
	decoder := yaml.NewDecoder(
		reader,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
		yaml.Strict(),
	)

	err = decoder.Decode(&spec)
	if err != nil {
		return spec, YamlParseError{
			Err:      fmt.Errorf("can`t parse yaml in '%s': %v", archFile, err),
			Warnings: []Warning{},
		}
	}

	specValidator := newValidator(spec, data, rootDirectory)
	warnings := specValidator.validate()
	if len(warnings) > 0 {
		return spec, YamlParseError{
			Err:      fmt.Errorf("spec '%s' has warnings", archFile),
			Warnings: warnings,
		}
	}

	return spec, YamlParseError{}
}
