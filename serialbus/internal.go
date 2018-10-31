package serialbus

import (
	"encoding/json"
)

func init() {
	emptyStruct := struct{}{}
	b, err := json.Marshal(emptyStruct)
	if err != nil {
		panic(err)
	}
	_FINISHSignal = string(b)
}
