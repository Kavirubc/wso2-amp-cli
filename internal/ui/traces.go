package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kavirubc/wso2-amp-cli/internal/api"
)

// Icons for traces
const (
	IconTrace = "📊"
)

// Tree drawing constants
const (
	TreeBranch     = "├─"
	TreeLastBranch = "└─"
	TreePipe       = "│ "
	TreeSpace      = "  "
)

// SpanNode for building span tree
type SpanNode struct {
	Span     api.Span
	Children []*SpanNode
}

// FormatNanosDuration converts nanoseconds to human-readable duration
func FormatNanosDuration(nanos int64) string {
	d := time.Duration(nanos)
	if d < time.Millisecond {
		return fmt.Sprintf("%dμs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fm", d.Minutes())
}

// TraceStatusCell returns a styled status based on error count
func TraceStatusCell(errorCount int) string {
	if errorCount == 0 {
		return SuccessStyle.Render("OK")
	}
	return ErrorStyle.Render(fmt.Sprintf("%d errors", errorCount))
}

// SpanStatusCell returns a styled span status
func SpanStatusCell(status string) string {
	switch strings.ToUpper(status) {
	case "OK":
		return SuccessStyle.Render("OK")
	case "ERROR":
		return ErrorStyle.Render("ERROR")
	default:
		return MutedStyle.Render(status)
	}
}

// TruncateTraceID shortens trace ID for table display
func TruncateTraceID(traceID string) string {
	if len(traceID) > 16 {
		return traceID[:16] + "..."
	}
	return traceID
}

// TruncateString shortens a string for display
func TruncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// BuildSpanTree builds a tree from flat span list using parentSpanId
func BuildSpanTree(spans []api.Span) []*SpanNode {
	// Build nodes map and find roots
	nodes := make(map[string]*SpanNode)
	var roots []*SpanNode

	for _, span := range spans {
		nodes[span.SpanID] = &SpanNode{Span: span}
	}

	for _, span := range spans {
		node := nodes[span.SpanID]
		if span.ParentSpanID == "" {
			roots = append(roots, node)
		} else if parent, ok := nodes[span.ParentSpanID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node) // orphan becomes root
		}
	}
	return roots
}

// RenderSpanTree renders the span tree with box-drawing chars
func RenderSpanTree(nodes []*SpanNode, prefix string, maxSpans int, currentCount *int) string {
	var sb strings.Builder
	for i, node := range nodes {
		if *currentCount >= maxSpans {
			remaining := countRemainingNodes(nodes[i:])
			if remaining > 0 {
				sb.WriteString(fmt.Sprintf("%s[+%d more spans...]\n", prefix, remaining))
			}
			break
		}
		*currentCount++

		isLast := i == len(nodes)-1
		branch := TreeBranch
		if isLast {
			branch = TreeLastBranch
		}

		duration := FormatNanosDuration(node.Span.DurationInNanos)
		status := SpanStatusCell(node.Span.Status)
		sb.WriteString(fmt.Sprintf("%s%s [%s] %s (%s)\n",
			prefix, branch, duration, node.Span.Name, status))

		childPrefix := prefix
		if isLast {
			childPrefix += TreeSpace
		} else {
			childPrefix += TreePipe
		}
		sb.WriteString(RenderSpanTree(node.Children, childPrefix, maxSpans, currentCount))
	}
	return sb.String()
}

// countRemainingNodes counts all nodes in slice including children
func countRemainingNodes(nodes []*SpanNode) int {
	count := 0
	for _, node := range nodes {
		count++ // count this node
		count += countRemainingNodes(node.Children)
	}
	return count
}

// FormatTraceTimestamp formats a trace timestamp for display
func FormatTraceTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
