// Package slicekit provides slice manipulation utilities.
//
// Includes:
//   - Filtering empty strings from slices
//   - Finding subsequences within slices
package slicekit

// RemoveEmpty filters out empty strings from a slice.
func RemoveEmpty(s []string) []string {
    var result []string
    for _, str := range s {
        if str != "" {
            result = append(result, str)
        }
    }
    return result
}


// ContainsSubSlice returns the index of the first occurrence of subSlice in slice.
// Returns -1 if subSlice is not found.
func ContainsSubSlice(slice, subSlice []string) int {
    // for i := 0; i <= len(slice)-len(subSlice); i++ {
    for i := range(len(slice)-len(subSlice)) {
        match := true
        // for j := 0; j < len(subSlice); j++ {
        for j := range(len(subSlice)) {
            if slice[i+j] != subSlice[j] {
                match = false
                break
            }
        }
        if match {
            return i
        }
    }
    return -1
}
