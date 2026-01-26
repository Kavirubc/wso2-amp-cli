package api

import (
	"net/url"
	"strconv"
)

// ListOptions contains pagination parameters for list operations
type ListOptions struct {
	Limit  int
	Offset int
}

// DefaultListOptions returns sensible defaults for pagination
func DefaultListOptions() ListOptions {
	return ListOptions{
		Limit:  10,
		Offset: 0,
	}
}

// buildPaginationQuery appends limit and offset to query params
func buildPaginationQuery(params url.Values, opts ListOptions) {
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
}
