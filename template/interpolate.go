package template

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	"github.com/cppforlife/go-patch/patch"
)

type InterpolateConfig struct {
	TemplateFile io.Reader
	VarsFiles    []io.Reader
	OpsFiles     []io.Reader
}

func (c *InterpolateConfig) Evaluate() ([]byte, error) {
	template, err := ioutil.ReadAll(c.TemplateFile)
	if err != nil {
		return []byte{}, err
	}

	tpl := boshtpl.NewTemplate(template)
	staticVars := boshtpl.StaticVariables{}
	ops := patch.Ops{}

	for _, file := range c.VarsFiles {
		var fileVars boshtpl.StaticVariables
		err = readYAMLFile(file, &fileVars)
		if err != nil {
			return nil, err
		}
		for k, v := range fileVars {
			staticVars[k] = v
		}
	}

	for _, file := range c.OpsFiles {
		var opDefs []patch.OpDefinition
		err = readYAMLFile(file, &opDefs)
		if err != nil {
			return nil, err
		}
		op, err := patch.NewOpsFromDefinitions(opDefs)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	evalOpts := boshtpl.EvaluateOpts{
		UnescapedMultiline: true,
		ExpectAllKeys:      true,
	}

	bytes, err := tpl.Evaluate(staticVars, ops, evalOpts)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func readYAMLFile(f io.Reader, dataType interface{}) error {
	payload, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(payload, dataType)
	if err != nil {
		return err
	}
	return nil
}
