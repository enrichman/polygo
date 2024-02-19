# Polygo

Polygo is a lightweight Go package for decoding polymorphic JSON responses effortlessly.

Dealing with APIs that return various types of objects can be challenging.  
Polygo simplifies this process by allowing you to map your types into a common interface.

## Example

Consider an API endpoint `/v1/shapes` that returns a list of _shapes_, each defined by a type field:

```json
[
    { "type": "circle", "radius": 5 },
    { "type": "square", "side":   3 }
]
```

With Polygo, you can easily handle this polymorphic JSON response. Here's how.

1. **Create a Decoder:** Initialize a decoder with a common interface and the field name used to check the object type.
2. **Register Types:** Register your concrete types with the decoder.
3. **Unmarshal JSON:** Use one of the available functions to unmarshal your JSON data.


```go
// Define your shape interface
type Shape interface {
    Area() float64
}

// Create a decoder specifying the interface and the field name,
// and register your concrete types
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

## Available functions

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

### UnmarshalInnerObject

`UnmarshalInnerObject` will unmarshal an object, looking into the specified path (using the [github.com/tidwall/gjson](github.com/tidwall/gjson) library).

```go
jsonBody := []byte(`{
    "data": { "type": "circle", "radius": 5 }
}`)

shapes, err := decoder.UnmarshalInnerObject("data", jsonBody)
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

## Installation
To use Polygo in your Go project, simply import it:

```go
import "github.com/enrichman/polygo"
```

## Contributing

Contributions are welcome! Feel free to open issues or pull requests on GitHub.

## License

[MIT](https://github.com/enrichman/polygo/blob/main/LICENSE)

## Feedback

If you like the project please star it on Github 🌟, and feel free to drop me a note on [Twitter](https://twitter.com/enrichmann)https://twitter.com/enrichmann, or open an issue!
