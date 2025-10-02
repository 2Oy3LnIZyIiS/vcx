package convertkit_test

import (
	"fmt"

	"vcx/pkg/toolkit/convertkit"
)

func Example() {
	// Convert boolean to integer
	fmt.Printf("true as int: %d\n", convertkit.BoolToInt(true))
	fmt.Printf("false as int: %d\n", convertkit.BoolToInt(false))

	// Convert integer to boolean
	fmt.Printf("5 as bool: %v\n", convertkit.IntToBool(5))
	fmt.Printf("0 as bool: %v\n", convertkit.IntToBool(0))
	fmt.Printf("-3 as bool: %v\n", convertkit.IntToBool(-3))

	// Output:
	// true as int: 1
	// false as int: 0
	// 5 as bool: true
	// 0 as bool: false
	// -3 as bool: false
}
