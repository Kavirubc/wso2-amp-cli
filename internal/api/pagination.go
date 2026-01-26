package api

import (
	"fmt"
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
// Negative values are treated as invalid and ignored
func buildPaginationQuery(params url.Values, opts ListOptions) {
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
}

// ValidatePaginationOptions validates that pagination options are non-negative
func ValidatePaginationOptions(opts ListOptions) error {
	if opts.Limit < 0 {
		return fmt.Errorf("limit must be non-negative, got %d", opts.Limit)
	}
	if opts.Offset < 0 {
		return fmt.Errorf("offset must be non-negative, got %d", opts.Offset)
	}
	return nil
}
