package utils

import (
	"fmt"
	"time"
)

// HumanDuration converts a time.Duration into a friendly string like "2 days ago"
func HumanDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	if days > 0 {
		return fmt.Sprintf("%d day%s ago", days, plural(days))
	}
	hours := int(d.Hours())
	if hours > 0 {
		return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
	}
	minutes := int(d.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%d minute%s ago", minutes, plural(minutes))
	}
	return "just now"
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
