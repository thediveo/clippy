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

package debug

import (
	"bytes"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("debug and tint logging", func() {

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
		SetWriter(rootCmd, &output)
		clippy.AddFlags(rootCmd)
	})

	It("defaults to logging infos or worse", func() {
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(BeEmpty())
		slog.Info("*information*")
		Expect(output.String()).To(MatchRegexp(`level=INFO msg=\*information\*`))
	})

	It("logs debugs", func() {
		rootCmd.SetArgs([]string{"foo", "--" + DebugFlagName})
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(MatchRegexp(`(?s)msg="debug logging enabled".*level=DEBUG msg=\*debug\*`))
	})

	It("changes the default", func() {
		rootCmd.SetArgs([]string{"foo"})
		SetDefaultLevel(rootCmd, slog.LevelDebug)
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Debug("*debug*")
		Expect(output.String()).To(MatchRegexp(`(?s)msg="debug logging enabled".*level=DEBUG msg=\*debug\*`))
	})

	It("resorts to stderr", func() {
		SetWriter(rootCmd, nil)
		Expect(rootCmd.Execute()).To(Succeed())
	})

	It("tints", func() {
		const (
			ansiBrightGreen = "\u001b\\[92m"
			ansiReset       = "\u001b\\[0m"
		)

		rootCmd.SetArgs([]string{"foo", "--" + TintedFlagName})
		Expect(rootCmd.Execute()).To(Succeed())
		slog.Info("hellorld!")
		Expect(output.String()).To(MatchRegexp(ansiBrightGreen + "INF" + ansiReset + " hellorld!"))
	})

	Context("created when necessary", func() {

		It("SetDefaultLevel", func() {
			cmd := &cobra.Command{}
			SetDefaultLevel(cmd, slog.LevelInfo)
			Expect(cmd.Context()).NotTo(BeNil())
		})

		It("SetLevel", func() {
			cmd := &cobra.Command{}
			SetLevel(cmd, slog.LevelInfo)
			Expect(cmd.Context()).NotTo(BeNil())
		})

	})

})
