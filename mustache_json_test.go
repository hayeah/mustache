package mustache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tailscale/hujson"
)

func TestRenderJSONTemplate(t *testing.T) {
	type UserData struct {
		Name string
		Age  int
	}

	// type BookRequest struct {
	// 	UserID string
	// 	BookID string
	// }

	tests := []struct {
		name        string
		template    string
		data        interface{}
		want        string
		expectedErr string
	}{
		{
			name:     "json object with struct data",
			template: `{"name": {{Name}}, "age": {{Age}}}`,
			data:     UserData{Name: "Alice", Age: 25},
			want:     `{"name":"Alice","age":25}`,
		},
		{
			name:     "json array of objects",
			template: `{"users": {{.}}}`,
			data: []UserData{
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 30},
			},
			want: `{"users":[{"Name":"Alice","Age":25},{"Name":"Bob","Age":30}]}`,
		},
		{
			name: "mustache section with slice of UserData",
			template: `{"users": [
			{{#.}} { "name": {{Name}} }, {{/.}}
]}`,
			data: []UserData{
				{Name: "Eve"},
				{Name: "Frank"},
			},
			want: `{"users":[{"name":"Eve"},{"name":"Frank"}]}`,
		},
		{
			name:     "missing variable in struct data",
			template: `{"name": {{Name}}, "height": {{Height}}}`,
			data:     UserData{Name: "Alice"}, // height is missing

			expectedErr: `missing variable "Height"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := RenderJSON(test.template, test.data)

			if test.expectedErr != "" {
				assert.Error(err)
				assert.EqualError(err, test.expectedErr)
				return
			}

			assert.NoError(err)

			out, err := hujson.Minimize(got)
			assert.NoError(err)
			assert.Equal(test.want, string(out))
		})
	}
}
