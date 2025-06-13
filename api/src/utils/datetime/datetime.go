package datetime

import (
	"time"

	"github.com/dromara/carbon/v2"
)

type Time     = carbon.Carbon
type Duration = time.Duration

// Now returns the current date and time in UTC.
//
// Returns:
//   - carbon.Carbon: The current date and time as a Carbon instance.
func Now() carbon.Carbon {
	return carbon.Now(carbon.UTC)
}

// Zero returns the zero date and time (January 1, 1 AD, 00:00:00) in UTC.
//
// Returns:
//   - carbon.Carbon: The zero date and time as a Carbon instance.
func Zero() carbon.Carbon {
	return carbon.Parse("")
}

// Parse parses a date and time string and returns a Carbon instance.
//
// Parameters:
//   - s: The date and time string to parse (e.g., "2024-12-14 15:04:05").
//
// Returns:
//   - carbon.Carbon: A Carbon instance representing the parsed date and time.
func Parse(s string) carbon.Carbon {
	return carbon.Parse(s)
}

// DiffInSeconds calculates the difference in seconds between two Carbon instances.
//
// Parameters:
//   - c1: The first Carbon instance.
//   - c2: The second Carbon instance.
//
// Returns:
//   - int: The difference in seconds between the two instances.
func DiffInSeconds(c1 carbon.Carbon, c2 carbon.Carbon) int {
	return int(c1.DiffInSeconds(c2))
}

// StartOfWeek returns the starting date and time of the current week, assuming
// the week starts on Monday.
//
// Returns:
//   - carbon.Carbon: A Carbon instance representing the start of the week.
func StartOfWeek() carbon.Carbon {
	return Now().SetWeekStartsAt(carbon.Monday).StartOfWeek()
}

// EndOfWeek returns the ending date and time of the current week, assuming
// the week ends on Sunday.
//
// Returns:
//   - carbon.Carbon: A Carbon instance representing the end of the week.
func EndOfWeek() carbon.Carbon {
	return Now().SetWeekStartsAt(carbon.Monday).EndOfWeek()
}

// FromTime converts a Go standard library `time.Time` instance into a Carbon instance.
//
// Parameters:
//   - t: The `time.Time` instance to convert.
//
// Returns:
//   - carbon.Carbon: A Carbon instance representing the given time.
func FromTime(t time.Time) carbon.Carbon {
	return carbon.CreateFromStdTime(t)
}

// ToTime converts a Carbon instance into a Go standard library `time.Time` instance.
//
// Parameters:
//   - c: The Carbon instance to convert.
//
// Returns:
//   - time.Time: A `time.Time` instance representing the given Carbon time.
func ToTime(c carbon.Carbon) time.Time {
	return c.StdTime()
}

// ToISO converts a Carbon instance into an ISO 8601 formatted string.
//
// Parameters:
//   - c: The Carbon instance to convert.
//
// Returns:
//   - string: The ISO 8601 formatted string representation of the date and time.
func ToISO(c carbon.Carbon) string {
	return c.ToIso8601String()
}
