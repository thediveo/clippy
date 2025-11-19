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
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/thediveo/clippy/cliplugin"
	"github.com/thediveo/go-plugger/v3"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func init() {
	plugger.Group[cliplugin.SetupCLI]().Register(canarySetupCLI)
	plugger.Group[cliplugin.BeforeCommand]().Register(canaryBeforeRun)
}

type canaryVariable int

const (
	setupCLICountRef canaryVariable = iota
	beforeRunRef
	beforeRunErr
)

func canarySetupCLI(cmd *cobra.Command) {
	setupCLICount := cmd.Context().Value(setupCLICountRef).(*int)
	*setupCLICount++
}

func canaryBeforeRun(cmd *cobra.Command) error {
	beforeRun := cmd.Context().Value(beforeRunRef).(*int)
	*beforeRun++
	beforeRunErr, _ := cmd.Context().Value(beforeRunErr).(error)
	return beforeRunErr
}

var _ = Describe("clippy", func() {

	It("calls AddFlags plugin method", func() {
		var count int
		ctx := context.WithValue(context.Background(),
			setupCLICountRef, &count)

		rootCmd := &cobra.Command{}
		rootCmd.SetArgs([]string{})
		rootCmd.SetContext(ctx)
		AddFlags(rootCmd)
		Expect(count).To(Equal(1))
	})

	It("calls BeforeCommand plugin method", func() {
		var count int
		ctx := context.WithValue(context.Background(),
			beforeRunRef, &count)

		ctx = context.WithValue(ctx,
			beforeRunErr, errors.New("fooerror"))

		rootCmd := &cobra.Command{}
		rootCmd.SetArgs([]string{})
		rootCmd.SetContext(ctx)
		Expect(BeforeCommand(rootCmd)).To(HaveOccurred())
		Expect(count).To(Equal(1))

		ctx = context.WithValue(ctx,
			beforeRunErr, error(nil))
		rootCmd.SetContext(ctx)
		Expect(BeforeCommand(rootCmd)).ToNot(HaveOccurred())
		Expect(count).To(Equal(2))
	})

})
