package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileContextFrom(t *testing.T) {
	tests := []struct {
		name        string
		giveSource  map[string]interface{}
		giveProfile string
		wantResult  profileContext
		wantErr     error
	}{
		{
			name: "should return invalid syntax when no 'contexts' key",
			giveSource: map[string]interface{}{
				"invalid": []interface{}{},
			},
			giveProfile: "",
			wantResult:  profileContext{},
			wantErr:     ErrInvalidContextsFileSyntax,
		},
		{
			name: "should return invalid syntax when different type of 'contexts' key",
			giveSource: map[string]interface{}{
				"invalid": "value",
			},
			giveProfile: "",
			wantResult:  profileContext{},
			wantErr:     ErrInvalidContextsFileSyntax,
		},
		{
			name: "should return profile not found when no entries in 'contexts' array",
			giveSource: map[string]interface{}{
				"contexts": []interface{}{},
			},
			giveProfile: "",
			wantResult:  profileContext{},
			wantErr:     ErrProfileNotFound,
		},
		{
			name: "should skip invalid syntax profiles",
			giveSource: map[string]interface{}{
				"contexts": []interface{}{
					map[string]interface{}{
						"invalid_key": 123,
					},
					map[interface{}]interface{}{
						"name_1": "default",
					},
					map[interface{}]interface{}{
						"name": "default",
					},
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url_1":      "something",
							"username_1": "something_2",
						},
					},
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url":        "something",
							"username_1": "something_2",
						},
					},
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url_1":    "something",
							"username": "something_2",
						},
					},
				},
			},
			giveProfile: "default",
			wantResult:  profileContext{},
			wantErr:     ErrProfileNotFound,
		},
		{
			name: "should return populated profile context",
			giveSource: map[string]interface{}{
				"contexts": []interface{}{
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url":      "something",
							"username": "something_2",
						},
					},
				},
			},
			giveProfile: "default",
			wantResult: profileContext{
				name: "default",
				http: httpSpec{
					url:      "something",
					username: "something_2",
				},
			},
			wantErr: nil,
		},
		{
			name: "should return first context in order when same name",
			giveSource: map[string]interface{}{
				"contexts": []interface{}{
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url":      "something",
							"username": "something_2",
						},
					},
					map[interface{}]interface{}{
						"name": "default",
						"http": map[interface{}]interface{}{
							"url":      "something_else",
							"username": "something_else_2",
						},
					},
				},
			},
			giveProfile: "default",
			wantResult: profileContext{
				name: "default",
				http: httpSpec{
					url:      "something",
					username: "something_2",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := profileContextFrom(tt.giveSource, tt.giveProfile)
			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
