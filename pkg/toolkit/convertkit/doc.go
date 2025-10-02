/*
Package convert provides utility functions for converting between different data types.

This package aims to simplify common type conversions in Go by providing a consistent
API for transforming values between different types. While many of these conversions
could be done with built-in Go expressions, using this package can improve code
readability and ensure consistent conversion behavior throughout an application.

Basic Usage:

	// Convert boolean to integer (1 for true, 0 for false)
	intValue := convert.BoolToInt(true) // intValue = 1

	// Convert integer to boolean (true for positive, false for zero/negative)
	boolValue := convert.IntToBool(5) // boolValue = true

The package currently focuses on simple conversions but may be expanded to include
more complex transformations in the future.
*/
package convertkit
