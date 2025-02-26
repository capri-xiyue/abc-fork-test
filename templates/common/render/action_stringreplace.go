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

package render

import (
	"context"
	"strings"

	spec "github.com/abcxyz/abc/templates/model/spec/v1beta4"
)

func actionStringReplace(ctx context.Context, sr *spec.StringReplace, sp *stepParams) error {
	var replacerArgs []string //nolint:prealloc // strings.NewReplacer has a weird input slice, it's less confusing to append rather than preallocate.
	for _, r := range sr.Replacements {
		toReplace, err := parseAndExecuteGoTmpl(r.ToReplace.Pos, r.ToReplace.Val, sp.scope)
		if err != nil {
			return err
		}
		replaceWith, err := parseAndExecuteGoTmpl(r.With.Pos, r.With.Val, sp.scope)
		if err != nil {
			return err
		}
		replacerArgs = append(replacerArgs, toReplace, replaceWith)
	}
	replacer := strings.NewReplacer(replacerArgs...)

	if err := walkAndModify(ctx, sp, sr.Paths, func(buf []byte) ([]byte, error) {
		return []byte(replacer.Replace(string(buf))), nil
	}); err != nil {
		return err
	}

	return nil
}
