package lint

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfigs_Search(t *testing.T) {
	tests := []struct {
		configs RuntimeConfigs
		path    string
		found   bool
	}{
		{nil, "a", false},
		{
			RuntimeConfigs{{IncludedPaths: []string{"a/*/*.proto"}}},
			"a/c/d.proto",
			true,
		},
		{
			RuntimeConfigs{{IncludedPaths: []string{"a/*.proto"}}},
			"ac/d.proto",
			false,
		},
		{
			RuntimeConfigs{{IncludedPaths: []string{"a/*.proto"}, ExcludedPaths: []string{"a/b*.proto"}}},
			"a/b.proto",
			false, // not found as the path is excluded.
		},
		{
			RuntimeConfigs{{IncludedPaths: []string{"a/*.proto"}, ExcludedPaths: []string{"a/b*.proto"}}},
			"a/c.proto",
			true,
		},
	}
	for _, test := range tests {
		_, err := test.configs.Search(test.path)
		if found := err == nil; found != test.found {
			t.Errorf("RuntimeConfigs.Search path %q returned error: %v, but expect found: %v", test.path, err, test.found)
		}
	}
}

func TestReadConfigsJSON(t *testing.T) {
	content := `
	[
		{
			"included_paths": ["a"],
			"excluded_paths": ["b"],
			"rule_configs": {
				"rule_a": {
					"status": "enabled",
					"category": "warning"
				}
			}
		}
	]
	`

	configs, err := ReadConfigsJSON(strings.NewReader(content))
	if err != nil {
		t.Errorf("ReadConfigsJSON returns error: %v", err)
	}

	expected := RuntimeConfigs{
		{
			IncludedPaths: []string{"a"},
			ExcludedPaths: []string{"b"},
			RuleConfigs: map[string]RuleConfig{
				"rule_a": {
					Status:   "enabled",
					Category: "warning",
				},
			},
		},
	}
	if !reflect.DeepEqual(configs, expected) {
		t.Errorf("ReadConfigsJSON returns %q, but want %q", configs, expected)
	}
}
