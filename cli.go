// Copyright 2020, 2025 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package clippy

import (
	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// AddFlags runs all registered [cliplugin.SetupCLI] plugin functions in order
// to register CLI flags for the specified root command.
func AddFlags(rootCmd *cobra.Command) {
	for _, setupCLI := range plugger.Group[cliplugin.SetupCLI]().Symbols() {
		setupCLI(rootCmd)
	}
}

// BeforeCommand runs all registered [cliplugin.BeforeCommand] plugin functions
// just before the selected command runs; it terminates as soon as the first
// plugin function returns a non-nil error.
func BeforeCommand(cmd *cobra.Command) error {
	for _, beforeCmd := range plugger.Group[cliplugin.BeforeCommand]().Symbols() {
		if err := beforeCmd(cmd); err != nil {
			return err
		}
	}
	return nil
}
