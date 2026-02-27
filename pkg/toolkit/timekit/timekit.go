// Package timekit provides time and date utilities.
//
// Includes:
//   - Unix epoch timestamps
//   - Fractional time with configurable precision (milliseconds to microseconds)
//   - RFC3339 formatted date/time strings
package timekit

import (
	"math"
	"time"
)

// GetEpoch returns the current Unix timestamp in seconds.
func GetEpoch() int64 {
	return time.Now().Unix()
}


// GetFractionalTime returns the fractional part of current time as an integer.
// Precision levels: 1=100ms, 2=10ms, 3=1ms, 4=100μs, 5=10μs, 6=1μs (default).
func GetFractionalTime(precision int) int {
    // Default to microsecond precision if 0
    if precision == 0 {
        precision = 6
    }

    // Get current time with nanosecond precision
    now := time.Now()

    // Extract just the fractional part (nanoseconds)
    nanos := now.Nanosecond()

    // Convert nanoseconds to the desired precision
    // 1 nanosecond = 10^-9 seconds
    // We want to convert to 10^-precision seconds
    scale := int(math.Pow10(9 - precision))

    // Scale the nanoseconds to desired precision and round
    result := (nanos + scale/2) / scale

    return result
}

// GetDateTime returns current date and time in RFC3339 format.
func GetDateTime() string {
    return time.Now().Format(time.RFC3339)
}
