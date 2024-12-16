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

package config

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"
)

type Config struct {
	DefaultDomain     string
	OpenwebUIImageTag string
	OllamaImageTag    string
}

/**
 * GetConfig retrieves the configuration from a Kubernetes config map.
 *
 * This function fetches the config map and extracts the necessary values to create a new Config object.
 * It returns the Config object or an error if any of the required values are missing or cannot be retrieved.
 */
func GetConfig() (*Config, error) {
	configMap, err := configMap()
	if err != nil {
		return nil, err
	}

	openwebUIImageTag, err := getConfigMapString(configMap, constants.OpenwebUIImageTag)
	if err != nil {
		return nil, err
	}

	ollamaImageTag, err := getConfigMapString(configMap, constants.OllamaImageTag)
	if err != nil {
		return nil, err
	}

	defaultDomain, err := getConfigMapString(configMap, constants.DefaultDomain)
	if err != nil {
		return nil, err
	}

	return &Config{
		DefaultDomain:     defaultDomain,
		OpenwebUIImageTag: openwebUIImageTag,
		OllamaImageTag:    ollamaImageTag,
	}, nil

}

/**
 * getConfigMapString retrieves a string value from a Kubernetes config map.
 *
 * This function takes a ConfigMap object and a key as input, and returns the corresponding string value.
 * If the key is not found in the ConfigMap's Data or BinaryData fields, it returns an error indicating that
 * the required key was not found.
 */
func getConfigMapString(configMap *corev1.ConfigMap, key string) (string, error) {
	if s, ok := configMap.Data[key]; ok {
		return s, nil
	}
	if b, ok := configMap.BinaryData[key]; ok {
		return string(b), nil
	}
	return "", fmt.Errorf("malformed Config Map: required key %q not found", key)
}

/**
 * configMap retrieves the configuration map from Kubernetes.
 *
 * This function fetches a ConfigMap object with the specified name and namespace,
 * and returns it or an error if any issues occur during retrieval.
 */
func configMap() (*corev1.ConfigMap, error) {
	ctx := context.Background()
	restConfig := ctrl.GetConfigOrDie()
	ctrlClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{}
	if err := ctrlClient.Get(ctx, types.NamespacedName{
		Name:      constants.AIChatWorspaceConfigMapName,
		Namespace: constants.AIChatWorkspaceNamespace,
	}, configMap); err != nil {
		return nil, err
	}

	return configMap, nil

}
