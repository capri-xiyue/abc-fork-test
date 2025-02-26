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

// Package render implements the template rendering related subcommands.
package render

// This file implements the "templates render" subcommand for installing a template.

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/abcxyz/abc/templates/common"
	"github.com/abcxyz/abc/templates/common/render"
	"github.com/abcxyz/abc/templates/common/templatesource"
	"github.com/abcxyz/pkg/cli"
)

type Command struct {
	cli.BaseCommand
	flags RenderFlags
	// used in prompt UT.
	skipPromptTTYCheck bool
}

// Desc implements cli.Command.
func (c *Command) Desc() string {
	return "instantiate a template to setup a new app or add config files"
}

// Help implements cli.Command.
func (c *Command) Help() string {
	return `
Usage: {{ COMMAND }} [options] <source>

The {{ COMMAND }} command renders the given template.

The "<source>" is the location of the template to be rendered. This may have a
few forms:

  - A remote GitHub or GitLab repo with either a version @tag or with the magic
    version "@latest". Examples:
    - github.com/abcxyz/abc/t/rest_server@latest
    - github.com/abcxyz/abc/t/rest_server@v0.3.1
  - A local directory, like /home/me/mydir
  - (Deprecated) A go-getter-style location, with or without ?ref=foo. Examples:
    - github.com/abcxyz/abc.git//t/react_template?ref=latest
	- github.com/abcxyz/abc.git//t/react_template
`
}

// Flags implements cli.Command.
func (c *Command) Flags() *cli.FlagSet {
	set := c.NewFlagSet()
	c.flags.Register(set)
	return set
}

func (c *Command) Run(ctx context.Context, args []string) error {
	if err := c.Flags().Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	fs := &common.RealFS{}
	if err := destOK(fs, c.flags.Dest); err != nil {
		return err
	}

	wd, err := c.WorkingDir()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	backupDir := filepath.Join(
		homeDir,
		".abc",
		"backups",
		fmt.Sprint(time.Now().Unix()))

	downloader, err := templatesource.ParseSource(ctx, &templatesource.ParseSourceParams{
		CWD:         wd,
		Source:      c.flags.Source,
		GitProtocol: c.flags.GitProtocol,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	return render.Render(ctx, &render.Params{ //nolint:wrapcheck
		BackupDir:            backupDir,
		Backups:              true,
		Clock:                clock.New(),
		Cwd:                  wd,
		DebugScratchContents: c.flags.DebugScratchContents,
		DebugStepDiffs:       c.flags.DebugStepDiffs,
		DestDir:              c.flags.Dest,
		Downloader:           downloader,
		ForceOverwrite:       c.flags.ForceOverwrite,
		FS:                   fs,
		GitProtocol:          c.flags.GitProtocol,
		KeepTempDirs:         c.flags.KeepTempDirs,
		Inputs:               c.flags.Inputs,
		InputFiles:           c.flags.InputFiles,
		Manifest:             c.flags.Manifest,
		Prompt:               c.flags.Prompt,
		Prompter:             c,
		SkipInputValidation:  c.flags.SkipInputValidation,
		SkipPromptTTYCheck:   c.skipPromptTTYCheck,
		SourceForMessages:    c.flags.Source,
		Stdout:               c.Stdout(),
	})
}

// destOK makes sure that the output directory looks sane.
func destOK(fs fs.StatFS, dest string) error {
	fi, err := fs.Stat(dest)
	if err != nil {
		if common.IsStatNotExistErr(err) {
			return nil
		}
		return fmt.Errorf("os.Stat(%s): %w", dest, err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("the destination %q exists but isn't a directory", dest)
	}

	return nil
}
