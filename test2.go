package main

import (
	"encoding/json"
	"fmt"
)

type S struct {
	A uint32 `json:"a"`
	B string `json:"b"`
	C uint32 `json:"c"`
}

type S1 struct {
	B string `json:"b"`
	C uint32 `json:"c"`
	D uint32 `json:"d"`
}

func main() {
	s := S{
		A: 12,
		C: 2,
	}
	s1 := S1{
		B: "123",
		C: 99999,
		D: 10,
	}
	js, _ := json.Marshal(s)
	js1, _ := json.Marshal(s1)
	fmt.Println(js)

	var m map[string]interface{}
	json.Unmarshal(js, &m)

	json.Unmarshal(js1, &m)

	res, _ := json.Marshal(m)

	fmt.Println(string(res)) // {"a":12,"b":"123","c":99999,"d":10}
}
