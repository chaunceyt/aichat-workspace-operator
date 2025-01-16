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

package modelfiles

import "fmt"

// GetSystemPromptPattern returns a system prompt pattern based on the provided model and pattern.
// It calls the internal prompt function to generate the pattern.
func GetSystemPromptPattern(model, pattern string, preference string) string {
	return prompt(model, pattern, preference)
}

// prompt generates a system prompt template with default parameters for temperature, top_p, and seed.
// The model and pattern are used as placeholders in the generated template.
//
// Args:
//
//	model (string): The name of the model to be used in the prompt template.
//	pattern (string): The pattern to be included in the system prompt template.
//
// Returns:
//
//	string: A formatted string representing the system prompt template with default parameters and the provided model and pattern.
//
// Ref: https://github.com/ollama/ollama/blob/main/docs/modelfile.md#parameter
func prompt(model, pattern string, preference string) string {
	var promptTemplate string

	var stableTemplate = `
FROM %s
	
PARAMETER temperature 0.1
PARAMETER top_p 0.25
PARAMETER seed 42

SYSTEM """
%s"""
		`

	var balancedTemplate = `
FROM %s

PARAMETER temperature 0.5
PARAMETER top_p 0.4
PARAMETER seed 21
	
SYSTEM """
%s"""	
		`

	var creativeTemplate = `
FROM %s

PARAMETER temperature 0.9
PARAMETER top_p 0.7
PARAMETER seed 0

SYSTEM """
%s"""
		`

	switch preference {
	case "stable":
		promptTemplate = stableTemplate
	case "balanced":
		promptTemplate = balancedTemplate
	case "creative":
		promptTemplate = creativeTemplate
	default:
		promptTemplate = stableTemplate
	}

	return fmt.Sprintf(promptTemplate, model, pattern)
}
