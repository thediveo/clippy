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

package debug

import (
	"cmp"
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"
)

// Names of the CLI flags defined and used in this package.
const (
	DebugFlagName  = "debug"
	TintedFlagName = "tinted"
)

// Register our plugin functions for delayed registration of CLI flags we bring
// into the game and the things to check or carry out before the selected
// command is finally run.
func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(
		setupCLI, plugger.WithPlugin("clippy/debug"))
	plugger.Group[cliplugin.BeforeCommand]().Register(
		beforeCommand, plugger.WithPlugin("clippy/debug"))
}

// setupCLI adds the "--debug" and "--tinted" flags to the specified command that
// changes the logging level to debug or enable logging at the info level.
func setupCLI(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(DebugFlagName, false, "enables debug structured logging output")
	cmd.PersistentFlags().Bool(TintedFlagName, false, "tints structured logging output")
}

type ctxIoWriterKey int

const ctxIoWriter ctxIoWriterKey = 0

// beforeCommand enables debug logging (and tinting) before any command finally
// is executed.
func beforeCommand(cmd *cobra.Command) error {
	w := cmp.Or[io.Writer](cmd.Context().Value(ctxIoWriter).(io.Writer), os.Stderr)
	level := slog.LevelInfo
	if debug, _ := cmd.PersistentFlags().GetBool(DebugFlagName); debug {
		level = slog.LevelDebug
	}
	var handler slog.Handler
	if tinted, _ := cmd.PersistentFlags().GetBool(TintedFlagName); tinted {
		handler = tint.NewHandler(w, &tint.Options{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: level,
		})
	}
	slog.SetDefault(slog.New(handler))
	slog.Debug("debug logging enabled")
	return nil
}
