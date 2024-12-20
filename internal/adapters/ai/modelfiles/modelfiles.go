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
func GetSystemPromptPattern(model, pattern string) string {
	return prompt(model, pattern)
}

// prompt generates a system prompt template with default parameters for temperature, top_p, top_k, and seed.
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
func prompt(model, pattern string) string {
	var promptTemplate = `
FROM %s
	
PARAMETER temperature 0.1
PARAMETER top_p 0.5
PARAMETER top_k 40
PARAMETER seed 1
	
SYSTEM """
%s"""	
		`

	return fmt.Sprintf(promptTemplate, model, pattern)
}
