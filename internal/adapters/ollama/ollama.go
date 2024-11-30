package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	ollama "github.com/ollama/ollama/api"
)

// https://github.com/ollama/ollama/blob/main/docs/api.md

// PullModel Download a model from the ollama library.
// https://github.com/ollama/ollama/blob/main/docs/api.md#pull-a-model
// TODO add support for huggingface
func PullModel(modelName string, defaultBaseURL string) error {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	req := &ollama.PullRequest{
		Model: modelName,
	}

	progressFunc := func(resp ollama.ProgressResponse) error {
		fmt.Printf("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
		return nil
	}

	err = client.Pull(ctx, req, progressFunc)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/ollama/ollama/blob/main/docs/api.md#list-local-models
func DoesModelExist(modelName string, defaultBaseURL string) (bool, error) {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return false, err
	}

	client := ollama.NewClient(baseClientURL, httpClient)
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	rp, err := client.List(ctx)
	if err != nil {
		return false, err
	}

	for _, llm := range rp.Models {
		if llm.Model == modelName {
			return true, nil
		}
	}

	return false, nil
}

// CreateFromModelFile
// https://github.com/ollama/ollama/blob/main/docs/api.md#create-a-model
func CreateFromModelFile(modelName, defaultBaseURL string) (bool, error) {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return false, err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	progressFunc := func(resp ollama.ProgressResponse) error {
		fmt.Printf("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
		return nil
	}

	ctx := context.Background()

	for p, filename := range modelfiles {
		createModelName := fmt.Sprintf("%s-%s", modelName, p)
		modelfile := fmt.Sprintf(filename, modelName)
		fmt.Printf("Creating %s from %s pattern modelfile\n", createModelName, p)
		err = client.Create(ctx, &ollama.CreateRequest{
			Model:     createModelName,
			Modelfile: modelfile,
		}, progressFunc)

		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func RunningModels(defaultBaseURL string) ([]string, error) {
	httpClient := http.DefaultClient

	var models []string

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return models, err
	}

	client := ollama.NewClient(baseClientURL, httpClient)
	if err != nil {
		return models, err
	}

	ctx := context.Background()

	rp, err := client.ListRunning(ctx)
	if err != nil {
		return models, err
	}

	for _, llm := range rp.Models {
		models = append(models, llm.Name)
	}

	return models, nil
}
