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
	"context"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thediveo/go-plugger/v3"
	"golang.org/x/sync/errgroup"
)

// Do (background) work, such as serving or collecting, in a plugin function
// while the passed context yet isn't done and return after the work is done or
// has been terminated. Please note that in case the work immediately fails it
// then should return a helpful error. If the work was terminated by the
// context, a Do plugin function should return a nil error instead in the spirit
// of being correctly wound down.
type Do func(context.Context, *cobra.Command) error

// DoAll starts all registered (background) work Do plugin functions and then
// waits for them to return. If any work function returns an error, the context
// passed to all Do work functions will be cancelled and the work functions
// should then return, but without error. Do then either returns the earliest
// reported error or nil in case all Do functions gracefully ended.
func DoAll(ctx context.Context, cmd *cobra.Command) error {
	name := cmd.DisplayName()
	defer slog.Info(name + " work ended")

	// The group context automatically gets cancelled when either our passed
	// context is cancelled or when one of the Do worker functions returns an
	// error. This automatically tells all remaining Do worker functions to
	// gracefully wind down their work. Cancelling the passed-in context is not
	// considered an error but the way of life, so Do worker functions should
	// never return any non-nil error under these circumstances.
	group, ctx := errgroup.WithContext(ctx)
	slog.Info(name + " work starting")
	for _, plug := range plugger.Group[Do]().PluginsSymbols() {
		group.Go(func() error { return plug.S(ctx, cmd) })
	}
	return group.Wait()
}
