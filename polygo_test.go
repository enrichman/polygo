package polygo_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/enrichman/polygo"
	"github.com/stretchr/testify/assert"
)

type Shape interface {
	GetType() string
	Area() float64
}

type Circle struct {
	Type   string  `json:"type"`
	Radius float64 `json:"radius"`
}

func NewCircle(radius float64) *Circle {
	return &Circle{Type: "circle", Radius: radius}
}

func (c *Circle) GetType() string {
	return c.Type
}

func (c *Circle) Area() float64 {
	return math.Phi * math.Pow(c.Radius, 2)
}

type Square struct {
	Type string  `json:"type"`
	Side float64 `json:"side"`
}

func NewSquare(side float64) *Square {
	return &Square{Type: "square", Side: side}
}

func (s *Square) GetType() string { return s.Type }

func (s *Square) Area() float64 {
	return math.Pow(s.Side, 2)
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
				"radius": 5
			}`),
			expectedObj: NewCircle(5),
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
				"radius": 5
			}`),
			expectedErr: "field 'type' not found",
		},
	}

	decoder := polygo.NewDecoder[Shape]("type").
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
			json: []byte(`[
				{ "type": "circle", "radius": 5 },
				{ "type": "square", "side": 3 }
			]`),
			expectedObj: []Shape{
				NewCircle(5),
				NewSquare(3),
			},
		},
	}

	decoder := polygo.NewDecoder[Shape]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shapes, err := decoder.UnmarshalArray(tc.json)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedObj, shapes)
		})
	}
}

func Test_UnmarshalInnerArray(t *testing.T) {
	tt := []struct {
		name        string
		json        []byte
		path        string
		expectedObj any
		expectedErr string
	}{
		{
			name: "simple data response",
			json: []byte(`{
				"data": [
					{ "type": "circle", "radius": 5 },
					{ "type": "square", "side": 3 }
				]
			}`),
			path: "data",
			expectedObj: []Shape{
				NewCircle(5),
				NewSquare(3),
			},
		},
	}

	decoder := polygo.NewDecoder[Shape]("type").
		Register("circle", Circle{}).
		Register("square", Square{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shape, err := decoder.UnmarshalInnerArray(tc.path, tc.json)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedObj, shape)
		})
	}
}

func Test_UnmarshalInnerArrayInResponse(t *testing.T) {
	type Response struct {
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}

	jsonData := []byte(`{
		"message": "response returned",
		"data": [
			{ "type": "circle", "radius": 5 },
			{ "type": "square", "side": 3 }
		]
	}`)

	decoder := polygo.NewDecoder[Shape]("type").
		Register("circle", &Circle{}).
		Register("square", Square{})

	var resp Response
	err := json.Unmarshal(jsonData, &resp)
	assert.NoError(t, err)
	assert.Equal(t, "response returned", resp.Message)

	shapes, err := decoder.UnmarshalArray(resp.Data)
	assert.NoError(t, err)

	expectedObj := []Shape{
		NewCircle(5),
		NewSquare(3),
	}
	assert.Equal(t, expectedObj, shapes)
}
