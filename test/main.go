package main

import (
	"encoding/json"
	"fmt"

	"github.com/enrichman/polygo"
)

type Wheeler interface {
	Wheels() int
}

type Vehicle struct {
	Type string `json:"type"`
}

func (v *Vehicle) Wheels() int {
	return 5
}

type Truck struct {
	Vehicle
}

type Car struct {
	Vehicle
}

type WheelerList []Wheeler

func main() {
	decoder := polygo.NewDecoder[Wheeler]("type").
		Register("truck", Truck{}).
		Register("car", Car{})

	type Resp struct {
		Data json.RawMessage `json:"data"`
	}
	var resp Resp
	err := json.Unmarshal([]byte(`{
			"data": [{
				"type": "truck",
				"name": "my truck"
			},{
				"type": "car",
				"name": "my truck"
			}]}`), &resp)
	fmt.Println(err)

	vehicleArr, err := decoder.UnmarshalArray(resp.Data)
	fmt.Printf("%#v - %+v\n", vehicleArr, err)

	vehicleArr, err = decoder.UnmarshalArray([]byte(`[{
			"type": "truck",
			"name": "my truck"
		},{
			"type": "car",
			"name": "my truck"
		}]`))
	fmt.Printf("%#v - %+v\n", vehicleArr, err)

	vehicleArr, err = decoder.UnmarshalArray([]byte(`[{
			"type": "truck",
			"name": "my truck"
		},{
			"type": "car",
			"name": "my truck"
		}]`))
	fmt.Printf("%#v - %+v\n", vehicleArr, err)

	vehicleArr, err = decoder.UnmarshalInnerArray("data", []byte(`{
		"data": [{
			"type": "truck",
			"name": "my truck"
		},{
			"type": "car",
			"name": "my truck"
		}]}`))
	fmt.Printf("%#v - %+v\n", vehicleArr, err)

	vehicle, err := decoder.UnmarshalObject([]byte(`{
		"type": "truck",
		"name": "my truck"
	}`))
	fmt.Println(vehicle, err)
	print(vehicle)

	vehicle, err = decoder.UnmarshalInnerObject("data", []byte(`{"data":{
		"type": "car",
		"name": "my car"
	}}`))
	fmt.Println(vehicle, err)
	print(vehicle)
}

func print(vehicle Wheeler) {
	switch t := vehicle.(type) {
	case *Truck:
		fmt.Println("t is truck", t.Type, t.Wheels())
	case *Car:
		fmt.Println("t is car", t.Type, t.Wheels())
	default:
		fmt.Println("t is wheeler (?)", t.Wheels())
	}
}
