package timekit

import (
	"math"
	"time"
)

// Get current epoch in seconds (returns int64)
func GetEpoch() int64 {
	return time.Now().Unix()
}


// Returns the current fractional time as an integer
// precision:
// 	0: defaults to 1 microsecond resolution
// 	1: 100 milliseconds
// 	2: 10 milliseconds
// 	3: 1 millisecond
// 	4: 100 microseconds
// 	5: 10 microseconds
// 	6: 1 microsecond
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

// Get current date and time in human readable format
func GetDateTime() string {
    return time.Now().Format(time.RFC3339)
}
