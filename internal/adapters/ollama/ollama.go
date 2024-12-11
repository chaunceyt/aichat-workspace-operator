/*
Copyright 2024 AIChatWorkspace Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	ollama "github.com/ollama/ollama/api"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/ai/modelfiles"
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
func CreateFromModelFile(modelName, defaultBaseURL string, patterns []string) (bool, error) {
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

	for _, pattern := range patterns {
		createModelName := fmt.Sprintf("%s-%s", modelName, pattern)
		modelfile := modelfiles.GetSystemPromptPattern(modelName, pattern)
		fmt.Printf("Creating %s from %s pattern modelfile\n", createModelName, pattern)
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
