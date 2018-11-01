package serialbus

type Request struct {
	id        string
	device    string
	operation string
	parameter string
}

type Reponse struct {
	id     string
	status string
	result string
}

var _FINISHSignal string
