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
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/go-cmp/cmp"

	"github.com/abcxyz/abc/templates/common"
	"github.com/abcxyz/abc/templates/common/templatesource"
	"github.com/abcxyz/pkg/testutil"
)

func TestWriteManifest(t *testing.T) {
	t.Parallel()

	clk := mockClock(t)

	cases := []struct {
		name             string
		dryRun           bool
		dlMeta           *templatesource.DownloadMetadata
		templateContents map[string]string
		destDirContents  map[string]string
		inputs           map[string]string
		outputHashes     map[string][]byte
		want             map[string]string
		wantErr          string
	}{
		{
			name: "simple_success_non_canonical",
			templateContents: map[string]string{
				"spec.yaml": "some stuff",
				"a.txt":     "some other stuff",
			},
			destDirContents: map[string]string{
				"a.txt": "some other stuff",
			},
			dlMeta: &templatesource.DownloadMetadata{
				IsCanonical: false,
			},
			inputs: map[string]string{
				"pizza":     "hawaiian",
				"pineapple": "deal with it",
			},
			outputHashes: map[string][]byte{
				"a.txt": []byte("fake_output_hash_32_bytes_sha256"),
			},
			want: map[string]string{
				"a.txt": "some other stuff",
				".abc/manifest_nolocation_2023-12-08T23:59:02.000000013Z.lock.yaml": `# Generated by the "abc templates" command. Do not modify.
api_version: cli.abcxyz.dev/v1beta3
kind: Manifest
creation_time: 2023-12-08T23:59:02.000000013Z
modification_time: 2023-12-08T23:59:02.000000013Z
template_location: ""
location_type: ""
template_version: ""
template_dirhash: h1:uh/nUYc3HpipWEon9kYOsvSrEadfu8Q9TdfBuHcnF3o=
inputs:
    - name: pineapple
      value: deal with it
    - name: pizza
      value: hawaiian
output_hashes:
    - file: a.txt
      hash: h1:ZmFrZV9vdXRwdXRfaGFzaF8zMl9ieXRlc19zaGEyNTY=
`,
			},
		},
		{
			name: "simple_success_canonical",
			templateContents: map[string]string{
				"spec.yaml": "some stuff",
				"a.txt":     "some other stuff",
			},
			destDirContents: map[string]string{
				"a.txt": "some other stuff",
			},
			dlMeta: &templatesource.DownloadMetadata{
				IsCanonical:     true,
				CanonicalSource: "github.com/foo/bar",
				LocationType:    "remote_git",
				HasVersion:      true,
				Version:         "v1.2.3",
			},
			inputs: map[string]string{
				"pizza":     "hawaiian",
				"pineapple": "deal with it",
			},
			outputHashes: map[string][]byte{
				"a.txt": []byte("fake_output_hash_32_bytes_sha256"),
			},
			want: map[string]string{
				"a.txt": "some other stuff",
				".abc/manifest_github.com%2Ffoo%2Fbar_2023-12-08T23:59:02.000000013Z.lock.yaml": `# Generated by the "abc templates" command. Do not modify.
api_version: cli.abcxyz.dev/v1beta3
kind: Manifest
creation_time: 2023-12-08T23:59:02.000000013Z
modification_time: 2023-12-08T23:59:02.000000013Z
template_location: github.com/foo/bar
location_type: remote_git
template_version: v1.2.3
template_dirhash: h1:uh/nUYc3HpipWEon9kYOsvSrEadfu8Q9TdfBuHcnF3o=
inputs:
    - name: pineapple
      value: deal with it
    - name: pizza
      value: hawaiian
output_hashes:
    - file: a.txt
      hash: h1:ZmFrZV9vdXRwdXRfaGFzaF8zMl9ieXRlc19zaGEyNTY=
`,
			},
		},
		{
			name: "dryrun_no_output",
			dlMeta: &templatesource.DownloadMetadata{
				IsCanonical:     false,
				CanonicalSource: "github.com/foo/bar",
			},
			dryRun: true,
			templateContents: map[string]string{
				"spec.yaml": "some stuff",
				"a.txt":     "some other stuff",
			},
			destDirContents: map[string]string{
				"a.txt": "some other stuff",
			},
			inputs: map[string]string{
				"pizza":     "hawaiian",
				"pineapple": "deal with it",
			},
			outputHashes: map[string][]byte{
				"a.txt": []byte("fake_output_hash_32_bytes_sha256"),
			},
			want: map[string]string{
				"a.txt": "some other stuff",
			},
		},
		{
			name: "no_inputs",
			templateContents: map[string]string{
				"spec.yaml": "some stuff",
				"a.txt":     "some other stuff",
			},
			destDirContents: map[string]string{
				"a.txt": "some other stuff",
			},
			dlMeta: &templatesource.DownloadMetadata{
				IsCanonical: false,
			},
			inputs: map[string]string{},
			outputHashes: map[string][]byte{
				"a.txt": []byte("fake_output_hash_32_bytes_sha256"),
			},
			want: map[string]string{
				"a.txt": "some other stuff",
				".abc/manifest_nolocation_2023-12-08T23:59:02.000000013Z.lock.yaml": `# Generated by the "abc templates" command. Do not modify.
api_version: cli.abcxyz.dev/v1beta3
kind: Manifest
creation_time: 2023-12-08T23:59:02.000000013Z
modification_time: 2023-12-08T23:59:02.000000013Z
template_location: ""
location_type: ""
template_version: ""
template_dirhash: h1:uh/nUYc3HpipWEon9kYOsvSrEadfu8Q9TdfBuHcnF3o=
inputs: []
output_hashes:
    - file: a.txt
      hash: h1:ZmFrZV9vdXRwdXRfaGFzaF8zMl9ieXRlc19zaGEyNTY=
`,
			},
		},
		{
			name: "no_outputs",
			templateContents: map[string]string{
				"spec.yaml": "some stuff",
				"a.txt":     "some other stuff",
			},
			destDirContents: map[string]string{},
			dlMeta: &templatesource.DownloadMetadata{
				IsCanonical: false,
			},
			inputs: map[string]string{
				"pizza":     "hawaiian",
				"pineapple": "deal with it",
			},
			outputHashes: map[string][]byte{},
			want: map[string]string{
				".abc/manifest_nolocation_2023-12-08T23:59:02.000000013Z.lock.yaml": `# Generated by the "abc templates" command. Do not modify.
api_version: cli.abcxyz.dev/v1beta3
kind: Manifest
creation_time: 2023-12-08T23:59:02.000000013Z
modification_time: 2023-12-08T23:59:02.000000013Z
template_location: ""
location_type: ""
template_version: ""
template_dirhash: h1:uh/nUYc3HpipWEon9kYOsvSrEadfu8Q9TdfBuHcnF3o=
inputs:
    - name: pineapple
      value: deal with it
    - name: pizza
      value: hawaiian
output_hashes: []
`,
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			templateDir := t.TempDir()
			destDir := t.TempDir()

			common.WriteAllDefaultMode(t, templateDir, tc.templateContents)
			common.WriteAllDefaultMode(t, destDir, tc.destDirContents)

			ctx := context.Background()
			err := writeManifest(ctx, &writeManifestParams{
				clock:        clk,
				destDir:      destDir,
				dlMeta:       tc.dlMeta,
				dryRun:       tc.dryRun,
				fs:           &common.RealFS{},
				inputs:       tc.inputs,
				outputHashes: tc.outputHashes,
				templateDir:  templateDir,
			})

			if diff := testutil.DiffErrString(err, tc.wantErr); diff != "" {
				t.Fatal(diff)
			}

			got := common.LoadDirWithoutMode(t, destDir)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("destination directory contents were not as expected (-got,+want): %s", diff)
			}
		})
	}
}

func mockClock(t *testing.T) *clock.Mock {
	t.Helper()

	clk := clock.NewMock()
	// We don't use UTC time here because we want to make sure local time
	// gets converted to UTC time before saving.
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		t.Fatalf("time.LoadLocation(): %v", err)
	}
	// This time has no particular significance.
	clk.Set(time.Date(2023, 12, 8, 15, 59, 2, 13, loc))
	return clk
}
