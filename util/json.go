package util

import (
	"encoding/json"
	"fmt"
)

func PrintObjInJson(in any) {
	byte, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Println(string(byte))
}
