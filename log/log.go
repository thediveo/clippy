// Copyright 2022, 2025 Harald Albrecht.
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

package log

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/clippy/debug"
	"github.com/thediveo/go-plugger/v3"
)

const (
	LogFlagName = "log"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		setupCLI, plugger.WithPlugin("clippy/log"), plugger.WithPlacement(">clippy/debug"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		beforeCommand, plugger.WithPlugin("clippy/log"), plugger.WithPlacement("<clippy/debug"))
}

// setupCLI runs after(!) the debug flag's setupCLI so that we can add our
// "--log" flag and make it mutually exclusive to the "--debug" flag.
func setupCLI(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(LogFlagName, false, "enables logging output")
	cmd.MarkFlagsMutuallyExclusive(LogFlagName, debug.DebugFlagName)
	debug.SetDefaultLevel(cmd, slog.LevelError)
}

// beforeCommand runs before(!) the debug flag's beforeCommand and lowers the
// logging bar when the "--log" flag has been specified with the command. It
// does so by attaching the forced level to the context of the command.
func beforeCommand(cmd *cobra.Command) error {
	if log, _ := cmd.PersistentFlags().GetBool(LogFlagName); log {
		debug.SetLevel(cmd, slog.LevelInfo)
	}
	return nil
}
