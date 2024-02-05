package polygo_test

import (
	"fmt"
	"testing"

	"github.com/enrichman/polygo"
)

type Type interface {
	GetType() string
}

type Shape interface {
	Area() float32
}

type TypedShape struct {
	Type  string `json:"type"`
	Color string `json:"color"`
}

func (t *TypedShape) GetType() string {
	return t.Type
}

type Circle struct {
	TypedShape
	Radius float32
}

type Square struct {
	TypedShape
	Side float32
}

func Test_UnmarshalObject(t *testing.T) {
	tt := []struct {
		name        string
		json        []byte
		expectedObj any
	}{
		{
			name: "",
			json: []byte(`{
				"type": "circle",
				"color": "red"
			}`),
			expectedObj: &Circle{
				TypedShape: TypedShape{
					Type:  "circle",
					Color: "red",
				},
				Radius: 5,
			},
		},
	}

	decoder := polygo.NewDecoder[Type]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shape, err := decoder.UnmarshalObject(tc.json)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(shape.GetType())
		})
	}
}
