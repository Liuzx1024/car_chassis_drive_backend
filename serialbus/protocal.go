package serialbus

type Request struct {
	id   string
	task struct {
		device    string
		operation string
		parameter string
	}
}

type Reponse struct {
	id     string
	status string
	result string
}
