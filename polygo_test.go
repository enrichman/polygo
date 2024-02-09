package polygo_test

import (
	"testing"

	"github.com/enrichman/polygo"
	"github.com/stretchr/testify/assert"
)

type Type interface {
	GetType() string
}

type Shape interface {
	Area() float32
}

type ColouredShape struct {
	Type  string `json:"type"`
	Color string `json:"color"`
}

func (t *ColouredShape) GetType() string {
	return t.Type
}

type Circle struct {
	ColouredShape
	Radius float32 `json:"radius"`
}

type Square struct {
	ColouredShape
	Side float32
}

func Test_UnmarshalObject(t *testing.T) {
	tt := []struct {
		name        string
		json        []byte
		expectedObj any
		expectedErr string
	}{
		{
			name: "simple circle",
			json: []byte(`{
				"type": "circle",
				"color": "red",
				"radius": 5
			}`),
			expectedObj: &Circle{
				ColouredShape: ColouredShape{
					Type:  "circle",
					Color: "red",
				},
				Radius: 5,
			},
		},
		{
			name: "unknown type",
			json: []byte(`{
				"type": "unknown",
				"color": "red",
				"foo": "bar"
			}`),
			expectedErr: "type 'unknown' not registered",
		},
		{
			name: "field not present",
			json: []byte(`{
				"no-type": "circle",
				"color": "red",
				"radius": 5
			}`),
			expectedErr: "field 'type' not found",
		},
	}

	decoder := polygo.NewDecoder[Type]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shape, err := decoder.UnmarshalObject(tc.json)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedObj, shape)
		})
	}
}

func Test_UnmarshalArray(t *testing.T) {
	tt := []struct {
		name        string
		json        []byte
		expectedObj any
		expectedErr string
	}{
		{
			name: "simple array",
			json: []byte(`[{
				"type": "circle",
				"color": "red",
				"radius": 5
			}]`),
			expectedObj: []Type{&Circle{
				ColouredShape: ColouredShape{
					Type:  "circle",
					Color: "red",
				},
				Radius: 5,
			}},
		},
	}

	decoder := polygo.NewDecoder[Type]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shape, err := decoder.UnmarshalArray(tc.json)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedObj, shape)
		})
	}
}
