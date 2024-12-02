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

func CopyModel(sourceName, destinationName string, defaultBaseURL string) error {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	req := &ollama.CopyRequest{
		Source:      sourceName,
		Destination: destinationName,
	}

	err = client.Copy(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func CreateModel(modelName, modelFile string, defaultBaseURL string) error {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	req := &ollama.CreateRequest{
		Model:     modelName,
		Modelfile: modelFile,
		// Stream: false,
	}

	progressFunc := func(resp ollama.ProgressResponse) error {
		fmt.Printf("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
		return nil
	}

	err = client.Create(ctx, req, progressFunc)
	if err != nil {
		return err
	}

	return nil
}

func DeleteModel(modelName, defaultBaseURL string) error {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	req := &ollama.DeleteRequest{
		Model: modelName,
	}

	err = client.Delete(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func ShowModel(modelName, defaultBaseURL string) (ollama.ModelDetails, error) {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return *&ollama.ModelDetails{}, err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	req := &ollama.ShowRequest{
		Model: modelName,
	}

	rp, err := client.Show(ctx, req)
	if err != nil {
		return *&ollama.ModelDetails{}, err
	}
	fmt.Println(rp.Details)

	return rp.Details, nil
}

func ListModels(defaultBaseURL string) ([]string, error) {
	httpClient := http.DefaultClient

	var models = []string{}

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return models, err
	}

	client := ollama.NewClient(baseClientURL, httpClient)

	ctx := context.Background()

	rp, err := client.List(ctx)
	if err != nil {
		return models, err
	}

	for _, llm := range rp.Models {
		models = append(models, llm.Model)
	}

	return models, nil
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

func ListRunningModels(defaultBaseURL string) ([]string, error) {
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

func setOllamaHost(workspace string) string {
	serviceName := fmt.Sprintf("%s-ollama", workspace)
	ollamaPort := int64(11434)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, workspace, ollamaPort)
	return ollamaServerURI
}
