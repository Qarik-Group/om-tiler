package configurator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
)

type interpolateConfig struct {
	TemplateFile  string
	OpsFiles      []string
	VarsFiles     []string
	Vars          map[string]interface{}
	TemplateStore http.FileSystem
}

func (c *interpolateConfig) evaluate() ([]byte, error) {
	template, err := c.readFile(c.TemplateFile)
	if err != nil {
		return []byte{}, err
	}

	tpl := boshtpl.NewTemplate(template)
	staticVars := boshtpl.StaticVariables{}
	ops := patch.Ops{}

	for _, file := range c.OpsFiles {
		var opDefs []patch.OpDefinition
		err = c.readYAMLFile(file, &opDefs)
		if err != nil {
			return nil, err
		}
		op, err := patch.NewOpsFromDefinitions(opDefs)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	for _, file := range c.VarsFiles {
		var fileVars boshtpl.StaticVariables
		err = c.readYAMLFile(file, &fileVars)
		if err != nil {
			return nil, err
		}
		for k, v := range fileVars {
			staticVars[k] = v
		}
	}

	for k, v := range c.Vars {
		staticVars[k] = v
	}

	evalOpts := boshtpl.EvaluateOpts{
		UnescapedMultiline: true,
		ExpectAllKeys:      true,
		ExpectAllVarsUsed:  true,
	}

	bytes, err := tpl.Evaluate(staticVars, ops, evalOpts)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c *interpolateConfig) readFile(file string) ([]byte, error) {
	if filepath.Ext(file) == "" {
		file = fmt.Sprintf("%s.yml", file)
	}
	f, err := c.TemplateStore.Open(file)
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadAll(f)
}

func (c *interpolateConfig) readYAMLFile(file string, dataType interface{}) error {
	payload, err := c.readFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(payload, dataType)
	if err != nil {
		return err
	}
	return nil
}
