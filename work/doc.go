/*
Package work provides running plugin work functions in parallel, such as for
continously servicing HTTP endpoints, running continuous background work, and
more. All registered work executes in parallel and should gracefully wind down
when asked nicely by cancelling the passed context.

Work functions can also terminate successfully early (such as in case of
one-shot work) and then should simply return nil.

If a work function encounters problems either early on or later and before the
passed context is cancelled/done, it should simply return an error. All
remaining work functions still standing will then see their passed context
cancelled and should more or less gracefully wind down without returning errors.
*/
package work
