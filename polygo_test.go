package polygo_test

import (
	"fmt"
	"testing"

	"github.com/enrichman/polygo"
)

type Shape struct {
	Type string `json:"type"`
	Area float32
}

type Circle struct {
	Shape
	Radius float32
}

type Square struct {
	Shape
	Side float32
}

func Test_UnmarshalArray(t *testing.T) {
	tt := []struct {
		name string
		json string
	}{
		{
			name: "",
			json: `[
				{ "type": "truck", "name": "my truck" },
				{ "type": "car", "name": "my truck" }
			]`,
		},
		{
			name: "non existing",
			json: `[
				{ "type": "truck", "name": "my truck" }
			]`,
		},
	}

	decoder := polygo.NewDecoder[Shape]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shapes, err := decoder.UnmarshalArray([]byte(tc.json))
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(shapes[0].Area)
		})
	}
}
