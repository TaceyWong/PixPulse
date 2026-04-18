package converter

type ConvertRequest struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	Mode       string `json:"mode"` // "color" | "bw" | "render"
}

type ConvertResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
