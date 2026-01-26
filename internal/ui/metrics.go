package ui

import (
	"fmt"
	"time"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
)

// Icons for metrics
const (
	IconMetrics = "\U0001F4CA" // 📊
)

// FormatCPUValue formats CPU value (fraction) as percentage or millicores
func FormatCPUValue(value float64) string {
	// Value is typically in cores (e.g., 0.1 = 100m = 10%)
	if value < 0.001 {
		return "0m"
	}
	// Display as millicores
	millicores := value * 1000
	if millicores < 1 {
		return fmt.Sprintf("%.2fm", millicores)
	}
	return fmt.Sprintf("%.0fm", millicores)
}

// FormatCPUPercentage formats CPU value as percentage
func FormatCPUPercentage(value float64) string {
	// Value is typically a fraction where 1.0 = 100%
	percentage := value * 100
	if percentage < 0.1 {
		return "0.0%"
	}
	return fmt.Sprintf("%.1f%%", percentage)
}

// FormatMemoryValue formats memory value in bytes to human-readable format
func FormatMemoryValue(bytes float64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
	)

	if bytes < 0 {
		return "0 B"
	}
	if bytes < KiB {
		return fmt.Sprintf("%.0f B", bytes)
	}
	if bytes < MiB {
		return fmt.Sprintf("%.1f KiB", bytes/KiB)
	}
	if bytes < GiB {
		return fmt.Sprintf("%.0f MiB", bytes/MiB)
	}
	return fmt.Sprintf("%.2f GiB", bytes/GiB)
}

// FormatMetricTimestamp formats an RFC3339 timestamp for display
func FormatMetricTimestamp(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		// Try RFC3339Nano
		t, err = time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return timestamp
		}
	}
	return t.Format("2006-01-02 15:04:05")
}

// BuildCPUMetricsTable builds table data for CPU metrics
func BuildCPUMetricsTable(usage, requests, limits []api.MetricDataPoint) ([]string, [][]string) {
	headers := []string{"TIME", "USAGE", "REQUEST", "LIMIT"}

	// Use the longest slice to determine row count
	maxLen := len(usage)
	if len(requests) > maxLen {
		maxLen = len(requests)
	}
	if len(limits) > maxLen {
		maxLen = len(limits)
	}

	if maxLen == 0 {
		return headers, nil
	}

	rows := make([][]string, maxLen)
	for i := 0; i < maxLen; i++ {
		row := make([]string, 4)

		// Get timestamp from usage if available, otherwise from others
		var timestamp string
		if i < len(usage) {
			timestamp = usage[i].Timestamp
		} else if i < len(requests) {
			timestamp = requests[i].Timestamp
		} else if i < len(limits) {
			timestamp = limits[i].Timestamp
		}
		row[0] = FormatMetricTimestamp(timestamp)

		// Usage as percentage
		if i < len(usage) {
			row[1] = FormatCPUPercentage(usage[i].Value)
		} else {
			row[1] = "-"
		}

		// Request as millicores
		if i < len(requests) {
			row[2] = FormatCPUValue(requests[i].Value)
		} else {
			row[2] = "-"
		}

		// Limit as millicores
		if i < len(limits) {
			row[3] = FormatCPUValue(limits[i].Value)
		} else {
			row[3] = "-"
		}

		rows[i] = row
	}

	return headers, rows
}

// BuildMemoryMetricsTable builds table data for memory metrics
func BuildMemoryMetricsTable(usage, requests, limits []api.MetricDataPoint) ([]string, [][]string) {
	headers := []string{"TIME", "USAGE", "REQUEST", "LIMIT"}

	// Use the longest slice to determine row count
	maxLen := len(usage)
	if len(requests) > maxLen {
		maxLen = len(requests)
	}
	if len(limits) > maxLen {
		maxLen = len(limits)
	}

	if maxLen == 0 {
		return headers, nil
	}

	rows := make([][]string, maxLen)
	for i := 0; i < maxLen; i++ {
		row := make([]string, 4)

		// Get timestamp from usage if available, otherwise from others
		var timestamp string
		if i < len(usage) {
			timestamp = usage[i].Timestamp
		} else if i < len(requests) {
			timestamp = requests[i].Timestamp
		} else if i < len(limits) {
			timestamp = limits[i].Timestamp
		}
		row[0] = FormatMetricTimestamp(timestamp)

		// Usage
		if i < len(usage) {
			row[1] = FormatMemoryValue(usage[i].Value)
		} else {
			row[1] = "-"
		}

		// Request
		if i < len(requests) {
			row[2] = FormatMemoryValue(requests[i].Value)
		} else {
			row[2] = "-"
		}

		// Limit
		if i < len(limits) {
			row[3] = FormatMemoryValue(limits[i].Value)
		} else {
			row[3] = "-"
		}

		rows[i] = row
	}

	return headers, rows
}

// HasMetricsData checks if the metrics response contains any data
func HasMetricsData(metrics *api.MetricsResponse) bool {
	if metrics == nil {
		return false
	}
	return len(metrics.CpuUsage) > 0 ||
		len(metrics.CpuRequests) > 0 ||
		len(metrics.CpuLimits) > 0 ||
		len(metrics.Memory) > 0 ||
		len(metrics.MemoryRequests) > 0 ||
		len(metrics.MemoryLimits) > 0
}
