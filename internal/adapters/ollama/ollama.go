package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	ollama "github.com/ollama/ollama/api"
)

func PullModel(modelName string, defaultBaseURL string) error {
	httpClient := http.DefaultClient

	baseClientURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return err
	}

	client := ollama.NewClient(baseClientURL, httpClient)
	if err != nil {
		return err
	}

	ctx := context.Background()

	req := &ollama.PullRequest{
		Model: modelName,
	}

	progressFunc := func(resp api.ProgressResponse) error {
		fmt.Printf("Progress: status=%v, total=%v, completed=%v\n", resp.Status, resp.Total, resp.Completed)
		return nil
	}

	err = client.Pull(ctx, req, progressFunc)
	if err != nil {
		return err
	}

	return nil
}

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
