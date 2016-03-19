// Package blaster provides interfaces for scheduling and sending a
// number of concurrent requests in a "blast". A handy line when working
// with interfaces in go is:
//	var _ <interface type> = (*<struct type>)(nil)
// this helps to ensure that the struct you wrote actually implements the
// interface.
package blaster
