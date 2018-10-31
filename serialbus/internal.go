package serialbus

import (
	"encoding/json"
	"fmt"
)

var _FINISHSignal string

func init() {
	emptyStruct := struct{}{}
	b, err := json.Marshal(emptyStruct)
	if err != nil {
		panic(err)
	}
	_FINISHSignal = string(b)
	fmt.Println(_FINISHSignal)
}
