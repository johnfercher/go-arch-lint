package container

import (
	"fmt"
	"io/ioutil"

	"github.com/fe3dback/go-arch-lint/internal/glue/pathresolver"
	"github.com/fe3dback/go-arch-lint/internal/glue/specassembler"
	"github.com/fe3dback/go-arch-lint/internal/glue/yamlannotationparser"
	"github.com/fe3dback/go-arch-lint/internal/glue/yamlresolver"
	"github.com/fe3dback/go-arch-lint/internal/glue/yamlspecprovider"
)

func (c *Container) provideSpecAssembler(projectDir, moduleName, archFilePath string) *specassembler.SpecAssembler {
	return specassembler.NewSpecAssembler(
		c.provideYamlSpecProvider(archFilePath),
		c.providePathResolver(),
		c.provideSourceCodeReferenceResolver(archFilePath),
		projectDir,
		moduleName,
	)
}

func (c *Container) provideYamlSpecProvider(archFilePath string) *yamlspecprovider.YamlSpecProvider {
	return yamlspecprovider.NewYamlSpecProvider(
		c.provideSourceCode(archFilePath),
	)
}

func (c *Container) providePathResolver() *pathresolver.PathResolver {
	return pathresolver.NewPathResolver()
}

func (c *Container) provideSourceCode(filePath string) []byte {
	sourceCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("failed to provide source code of archfile: %s", err))
	}

	return sourceCode
}

func (c *Container) provideSourceCodeReferenceResolver(filePath string) *yamlresolver.YamlReferenceResolver {
	return yamlresolver.NewYamlReferenceResolver(
		c.provideSourceCode(filePath),
		filePath,
		c.provideAnnotationParser(),
	)
}

func (c *Container) provideAnnotationParser() *yamlannotationparser.AnnotationParser {
	return yamlannotationparser.NewAnnotationParser()
}