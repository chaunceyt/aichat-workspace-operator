package ollama

// modelfile
// https://github.com/ollama/ollama/blob/main/docs/modelfile.md
var modelfiles = map[string]string{
	"ai":            aiSystem,
	"createSummary": createSummary,
	"explainCode":   explainCode,
	"translate":     translate,
}
