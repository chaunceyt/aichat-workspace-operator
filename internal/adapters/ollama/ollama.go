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
/**
 * Downloads a model from the ollama library.
 *
 * @param modelName The name of the model to download.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return An error if the download fails, or nil otherwise.
 *
 * https://github.com/ollama/ollama/blob/main/docs/api.md#pull-a-model
 * TODO: add support for huggingface (confirm it works)
 */
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

/**
 * Copies a model from the source name to the destination name.
 *
 * @param sourceName The source name of the model to copy.
 * @param destinationName The destination name for the copied model.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return An error if the copy operation fails, or nil otherwise.
 */
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

/**
 * Creates a new model in the AIChat Workspace.
 *
 * @param modelName The name of the model to create.
 * @param modelFile The file path of the model to upload.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return An error if the creation fails, or nil otherwise.
 */
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

/**
 * Deletes a model from the AIChat Workspace.
 *
 * @param modelName The name of the model to delete.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return An error if the deletion fails, or nil otherwise.
 */
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

/**
 * Shows details of a model in the AIChat Workspace.
 *
 * @param modelName The name of the model to show details for.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return A ModelDetails object containing information about the model, or an error if the operation fails.
 */
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

/**
 * Lists all models in the AIChat Workspace.
 *
 * @param defaultBaseURL The base URL of the ollama API.
 * @return A list of model names as strings, or an error if the operation fails.
 */
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

/**
 * Checks if a model exists in the AIChat Workspace.
 *
 * @param modelName The name of the model to check for existence.
 * @param defaultBaseURL The base URL of the ollama API.
 * @return A boolean indicating whether the model exists, and an error if the operation fails.
 */
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

/**
 * Creates new models from the provided model name and SYSTEM prompt patterns.
 *
 * @param modelName The base model name to create.
 * @param defaultBaseURL The base URL of the ollama API.
 * @param patterns List of string patterns for creating multiple models with a single request.
 * @return A boolean indicating whether all creations were successful, or an error if any creation fails.
 *
 * https://github.com/ollama/ollama/blob/main/docs/api.md#create-a-model
 * Uses patterns from https://github.com/danielmiessler/fabric/tree/main/patterns
 * TODO: rename to CreateFromSystemPromptPattern
 */
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

/**
 * Lists the names of all running models in the AIChat Workspace.
 *
 * @param defaultBaseURL The base URL of the ollama API.
 * @return A list of model names as strings, or an error if the operation fails.
 */
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
