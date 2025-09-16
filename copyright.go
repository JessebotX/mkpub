package pub

type Copyright struct {
	Notice   string    `json:"notice"`
	Licenses []License `json:"licenses"`
}

type License struct {
	FileName string `json:"file_name"`
	Text     string `json:"text"`
}
