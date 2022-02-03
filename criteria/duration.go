package criteria

import (
	"strconv"
	"time"
)

func formatDuration(d time.Duration) string {
	days := uint64(d / (time.Hour * 24))
	d = d % (time.Hour * 24)

	hours := uint64(d / time.Hour)
	d = d % time.Hour

	mins := uint64(d / time.Minute)
	d = d % time.Minute

	secs := uint64(d / time.Second)

	s := ""

	if days > 0 {
		s = strconv.FormatUint(days, 10) + "d"
	}

	if hours > 0 {
		s += strconv.FormatUint(hours, 10) + "h"
	}

	if mins > 0 {
		s += strconv.FormatUint(mins, 10) + "m"
	}

	if secs > 0 {
		s += strconv.FormatUint(secs, 10) + "s"
	}

	return s
}
