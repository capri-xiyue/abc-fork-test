// Copyright 2023 The Authors (see AUTHORS file)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goldentest

import (
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/abcxyz/abc/templates/model"
)

// This file parses a YAML file that describes test configs.

// VarValue represents one of the parsed "input" fields from the inputs.yaml file.
type VarValue struct {
	// Pos is the YAML file location where this object started.
	Pos model.ConfigPos `yaml:"-"`

	Name  model.String `yaml:"name"`
	Value model.String `yaml:"value"`
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (i *VarValue) UnmarshalYAML(n *yaml.Node) error {
	return model.UnmarshalPlain(n, i, &i.Pos) //nolint:wrapcheck
}

func (i *VarValue) Validate() error {
	return errors.Join(
		model.NotZeroModel(&i.Pos, i.Name, "name"),
		model.NotZeroModel(&i.Pos, i.Value, "value"),
	)
}

// Test represents a parsed test.yaml describing test configs.
type Test struct {
	// Pos is the YAML file location where this object started.
	Pos model.ConfigPos `yaml:"-"`

	Inputs []*VarValue `yaml:"inputs"`
}

// Validate implements model.Validator.
func (t *Test) Validate() error {
	return errors.Join(
		model.ValidateEach(t.Inputs),
	)
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (t *Test) UnmarshalYAML(n *yaml.Node) error {
	return model.UnmarshalPlain(n, t, &t.Pos, "api_version", "apiVersion", "kind") //nolint:wrapcheck
}
