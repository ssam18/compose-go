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
	"fmt"

	"github.com/compose-spec/compose-go/v2/tree"
)

// transformImage converts the long syntax for an image reference to short syntax.
// Long syntax: {name: "myimage", tag: "1.0.0"} -> "myimage:1.0.0"
// Long syntax: {name: "myimage", digest: "sha256:abc"} -> "myimage@sha256:abc"
// Short syntax is returned unchanged.
func transformImage(data any, p tree.Path, _ bool) (any, error) {
	switch v := data.(type) {
	case string:
		return v, nil
	case map[string]any:
		return imageFromMap(v, p)
	default:
		return data, fmt.Errorf("%s: invalid type %T for image", p, data)
	}
}

// imageFromMap converts a map with name/tag/digest fields to an image reference string.
func imageFromMap(m map[string]any, p tree.Path) (string, error) {
	name, ok := m["name"].(string)
	if !ok || name == "" {
		return "", fmt.Errorf("%s: image long syntax requires a non-empty 'name' field", p)
	}

	digest, hasDigest := m["digest"].(string)
	tag, hasTag := m["tag"].(string)

	switch {
	case hasDigest && digest != "":
		return name + "@" + digest, nil
	case hasTag && tag != "":
		return name + ":" + tag, nil
	default:
		return name, nil
	}
}

// transformBuildTags converts a list of image references (short or long syntax) to a list of strings.
func transformBuildTags(data any, p tree.Path, ignoreParseError bool) (any, error) {
	entries, ok := data.([]any)
	if !ok {
		return data, fmt.Errorf("%s: invalid type %T for tags", p, data)
	}
	tags := make([]any, len(entries))
	for i, entry := range entries {
		t, err := transformImage(entry, p.Next("[]"), ignoreParseError)
		if err != nil {
			return nil, err
		}
		tags[i] = t
	}
	return tags, nil
}
