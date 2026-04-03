/*
   Copyright 2020 The Compose Specification Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package transform

import (
	"testing"

	"github.com/compose-spec/compose-go/v2/tree"
	"gotest.tools/v3/assert"
)

func Test_transformImage(t *testing.T) {
	p := tree.NewPath("services.foo.image")
	tests := []struct {
		name    string
		data    any
		want    any
		wantErr string
	}{
		{
			name: "string unchanged",
			data: "myimage:1.0.0",
			want: "myimage:1.0.0",
		},
		{
			name: "long syntax name and tag",
			data: map[string]any{"name": "myimage", "tag": "1.0.0"},
			want: "myimage:1.0.0",
		},
		{
			name: "long syntax name only",
			data: map[string]any{"name": "myimage"},
			want: "myimage",
		},
		{
			name: "long syntax name and digest",
			data: map[string]any{"name": "myimage", "digest": "sha256:abc123"},
			want: "myimage@sha256:abc123",
		},
		{
			name: "long syntax digest takes precedence over tag",
			data: map[string]any{"name": "myimage", "tag": "1.0.0", "digest": "sha256:abc123"},
			want: "myimage@sha256:abc123",
		},
		{
			name:    "map missing name",
			data:    map[string]any{"tag": "1.0.0"},
			wantErr: "services.foo.image: image long syntax requires a non-empty 'name' field",
		},
		{
			name:    "invalid type",
			data:    42,
			wantErr: "services.foo.image: invalid type int for image",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformImage(tt.data, p, false)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_transformBuildTags(t *testing.T) {
	p := tree.NewPath("services.foo.build.tags")
	tests := []struct {
		name    string
		data    any
		want    any
		wantErr string
	}{
		{
			name: "string list unchanged",
			data: []any{"myimage:1.0.0", "myimage:latest"},
			want: []any{"myimage:1.0.0", "myimage:latest"},
		},
		{
			name: "long syntax in list",
			data: []any{
				map[string]any{"name": "myimage", "tag": "1.0.0"},
				map[string]any{"name": "myimage", "tag": "latest"},
			},
			want: []any{"myimage:1.0.0", "myimage:latest"},
		},
		{
			name: "mixed short and long syntax",
			data: []any{
				map[string]any{"name": "myimage", "tag": "1.0.0"},
				"myimage:latest",
			},
			want: []any{"myimage:1.0.0", "myimage:latest"},
		},
		{
			name:    "not a list",
			data:    "myimage:1.0.0",
			wantErr: "services.foo.build.tags: invalid type string for tags",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformBuildTags(tt.data, p, false)
			if tt.wantErr != "" {
				assert.Error(t, err, tt.wantErr)
				return
			}
			assert.NilError(t, err)
			assert.DeepEqual(t, got, tt.want)
		})
	}
}
