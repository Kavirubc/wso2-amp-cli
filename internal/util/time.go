package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseSinceDuration parses strings like "1h", "24h", "7d" and returns the start time
func ParseSinceDuration(since string) (time.Time, error) {
	since = strings.TrimSpace(strings.ToLower(since))

	// Handle day suffix specially
	if strings.HasSuffix(since, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(since, "d"))
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid duration: %s", since)
		}
		if days <= 0 {
			return time.Time{}, fmt.Errorf("duration must be positive: %s", since)
		}
		return time.Now().Add(-time.Duration(days) * 24 * time.Hour), nil
	}

	// Use standard Go duration parsing for h, m, s
	duration, err := time.ParseDuration(since)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid duration: %s (use format like 1h, 30m, 24h, 7d)", since)
	}

	if duration <= 0 {
		return time.Time{}, fmt.Errorf("duration must be positive: %s", since)
	}

	return time.Now().Add(-duration), nil
}
