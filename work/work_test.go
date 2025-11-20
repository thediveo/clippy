// Copyright 2025 Harald Albrecht.
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

package work

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
)

var _ = Describe("service and background workers", func() {

	var output bytes.Buffer
	var group = plugger.Group[Do]()

	BeforeEach(func() {
		old := slog.Default()
		DeferCleanup(func() { slog.SetDefault(old) })

		output.Reset()
		slog.SetDefault(slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelInfo})))

		oldDoers := group.Backup()
		DeferCleanup(func() { group.Restore(oldDoers) })

		goodgos := Goroutines()
		DeferCleanup(func() {
			Eventually(Goroutines).Within(2 * time.Second).ProbeEvery(100 * time.Millisecond).
				ShouldNot(HaveLeaked(goodgos))
		})
	})

	It("informs when starting and ending work", func() {
		cmd := &cobra.Command{
			Use: "foo does bar",
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // sic!
		Expect(DoAll(ctx, cmd)).To(Succeed())
		Expect(output.String()).To(MatchRegexp(`(?s)level=INFO msg="foo work starting".*level=INFO msg="foo work ended"`))
	})

	It("ends all work gracefully", func(ctx context.Context) {
		cmd := &cobra.Command{
			Use: "foo does bar",
		}

		started := make(chan struct{}, 2)

		for _, no := range []string{"first", "second"} {
			group.Register(func(ctx context.Context, c *cobra.Command) error {
				defer GinkgoRecover()
				By("doing " + no + " work")
				started <- struct{}{}
				Eventually(ctx.Done).Within(5 * time.Second).Should(BeClosed())
				return nil
			})
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			defer GinkgoRecover()
			for range 2 {
				Eventually(started).WithContext(ctx).Within(3 * time.Second).Should(Receive())
			}
			cancel()
		}()

		Eventually(func() error { return DoAll(ctx, cmd) }).WithContext(ctx).
			Within(10 * time.Second).Should(Succeed())
	})

	It("ends remaining work upon failure", func() {
		cmd := &cobra.Command{
			Use: "foo does bar",
		}

		started := make(chan struct{})

		group.Register(func(ctx context.Context, c *cobra.Command) error {
			defer GinkgoRecover()
			By("doing first work initially and then failing")
			Eventually(started).WithContext(ctx).Within(5 * time.Second).Should(BeClosed())
			return errors.New("foo!")
		})

		group.Register(func(ctx context.Context, c *cobra.Command) error {
			defer GinkgoRecover()
			By("doing second work")
			close(started)
			Eventually(ctx.Done).Within(5 * time.Second).Should(BeClosed())
			return nil
		})

		Eventually(func() error { return DoAll(context.Background(), cmd) }).
			Within(10 * time.Second).Should(MatchError(ContainSubstring("foo!")))
	})

})
