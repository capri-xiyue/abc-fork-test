// Copyright 2023 The Authors (see AUTHORS file)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package common contains the common utility functions for template commands.

package templatesource

import (
	"context"
)

// A Downloader is returned by a sourceParser. It offers the ability to
// download a template, and provides some metadata.
type Downloader interface {
	// Download downloads this template into the given directory.
	Download(ctx context.Context, cwd, destDir string) (*DownloadMetadata, error)
}

type DownloadMetadata struct {
	// A "canonical" location is one that's the same for everybody. When
	// installing a template source like
	// "~/my_downloaded_templates/foo_template", that location is not canonical,
	// because not every has that directory downloaded on their machine. On the
	// other hand, a template location like
	// github.com/abcxyz/gcp-org-terraform-template *is* canonical because
	// everyone everywhere can access it by that name.
	//
	// Canonical template locations are preferred because they make automatic
	// template upgrades easier. Given a destination directory that is the
	// output of a template, we can easily upgrade it if we know the canonical
	// location of the template that created it. We just go look for new git
	// tags at the canonical location.
	//
	// A local template directory is not a canonical location except for one
	// special case: when the template source directory and the destination
	// directory are within the same repo. This supports the case where a single
	// git repo contains templates that are rendered into that repo. Since the
	// relative path between the template directory and the destination
	// directory are the same for everyone who clones the repo, that means the
	// relative path counts as a canonical source.
	//
	// IsCanonical is true if and only if CanonicalSource and LocationType are
	// non-empty.
	IsCanonical     bool
	CanonicalSource string
	LocationType    string

	// Depending on where the template was taken from, there might be a version
	// string associated with it (e.g. a git tag or a git SHA).
	//
	// HasVersion is true if and only if Version is non-empty.
	HasVersion bool
	Version    string

	// Values for template variables like _git_tag and _git_sha.
	Vars DownloaderVars
}

// Values for template variables like _git_tag and _git_sha.
type DownloaderVars struct {
	GitTag      string
	GitSHA      string
	GitShortSHA string
}
