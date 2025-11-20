// Copyright 2023, 2025 Harald Albrecht.
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
	"bytes"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"
	"github.com/thediveo/clippy/debug"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("info+ logging", func() {

	var rootCmd *cobra.Command
	var output bytes.Buffer

	BeforeEach(func() {
		oldLogger := slog.Default()
		DeferCleanup(func() {
			slog.SetDefault(oldLogger)
		})

		output.Reset()

		rootCmd = &cobra.Command{
			PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
				return clippy.BeforeCommand(cmd)
			},
			RunE: func(*cobra.Command, []string) error { return nil },
		}
		debug.SetWriter(rootCmd, &output)
		clippy.AddFlags(rootCmd)
	})

	It("defaults to logging only errors or worse", func() {
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(BeEmpty())
		slog.Info("*information*")
		Expect(output.String()).To(BeEmpty())
		slog.Error("*errør*")
		Expect(output.String()).To(MatchRegexp(`level=ERROR msg=\*errør\*`))
	})

	It("logs information when requested", func() {
		rootCmd.SetArgs([]string{"foo", "--" + LogFlagName})
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(BeEmpty())
		slog.Info("*inførmatiøn*")
		Expect(output.String()).To(MatchRegexp(`level=INFO msg=\*inførmatiøn\*`))
	})

	It("still logs debugs when requested", func() {
		rootCmd.SetArgs([]string{"foo", "--" + debug.DebugFlagName})
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(MatchRegexp(`(?s)msg="debug logging enabled".*level=DEBUG msg=\*debug\*`))
	})

})
