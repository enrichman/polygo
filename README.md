# Polygo

Decode polymorphic JSON response in a breeze with Go!

## Example

Sometimes you have to deal with APIs that return different kind of objects, usually distinguished by a field.  

In the following example the `/v1/shapes` endpoint returns a list of _shapes_, with the `type` defining the type of object:

```json
[
	{ "type": "circle", "radius": 5 },
	{ "type": "square", "side": 3 }
]
```

There are different ways on how to handle this, but with `polygo` you can just map your types in a common interface.

Create a decoder specifying the interface and the field name, register the concrete types, and unmarshal the JSON:
```go
// register your concrete types defined in a common interface
decoder := polygo.NewDecoder[Shape]("type").
    Register("circle", Circle{}).
    Register("square", Square{})

// unmarshal your JSON
shapes, _ := decoder.UnmarshalArray(jsonBytes)

for _, shape := range shapes {
    // use the methods defined by the interface
    fmt.Printf("shape area: %v\n", shape.Area())

    // or check the concrete type if needed
    switch s := shape.(type) {
    case *Circle:
        fmt.Printf("circle radius: %v\n", s.Radius)
    case *Square:
        fmt.Printf("square side: %v\n", s.Side)
    }
}
```

## Usage

Create a decoder with a common interface and the field name used to check the object type.   You can then register your types:


```go
type Shape interface {
    Area() float64
}

decoder := polygo.NewDecoder[Shape]("type").
    Register("circle", Circle{}).
    Register("square", Square{})
```

Then unmarshal your JSON with one of the available functions.

### UnmarshalObject

`UnmarshalObject` will unmarshal a plain object:

```go
jsonBody := []byte(`{ "type": "circle", "radius": 5 }`)

shape, err := decoder.UnmarshalObject(jsonBody)
if err != nil {
    return err
}
```

### UnmarshalArray

`UnmarshalArray` will unmarshal an array of objects:

```go
jsonBody := []byte(`[
    { "type": "circle", "radius": 5 },
    { "type": "square", "side": 3 }
]`)

shapes, err := decoder.UnmarshalArray(jsonBody)
if err != nil {
    return err
}
```

### UnmarshalInnerObjbect

`UnmarshalInnerObjbect` will unmarshal an object, looking into the specified path (using the [github.com/tidwall/gjson](github.com/tidwall/gjson) library).

```go
jsonBody := []byte(`{
    "data": { "type": "circle", "radius": 5 }
}`)

shapes, err := decoder.UnmarshalInnerObjbect("data", jsonBody)
if err != nil {
    return err
}
```

### UnmarshalInnerArray

`UnmarshalInnerArray` will unmarshal an array of objects, looking into the specified path (using the [github.com/tidwall/gjson](github.com/tidwall/gjson) library).

```go
jsonBody := []byte(`{
    "data": [
        { "type": "circle", "radius": 5 },
        { "type": "square", "side": 3 }
    ]
}`)

shapes, err := decoder.UnmarshalInnerArray("data", jsonBody)
if err != nil {
    return err
}
```

### Wrapped response

If your data is wrapped in an object with fields that you are interested to check, you should use a struct with a `json.RawMessage` field. Then you can unmarshal this field with the decoder.



```go
type Response struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    json.RawMessage `json:"data"`
}

jsonData := []byte(`{
    "code": 200,
    "message": "all good",
    "data": [
        { "type": "circle", "radius": 5 },
        { "type": "square", "side": 3 }
    ]
}`)

var resp Response
err := json.Unmarshal(jsonData, &resp)
if err != nil {
    return err
}

shapes, err := decoder.UnmarshalArray(resp.Data)
if err != nil {
    return err
}
```
