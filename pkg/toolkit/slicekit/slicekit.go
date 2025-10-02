package slicekit

// removeEmpty filters out empty strings from a slice of strings
func RemoveEmpty(s []string) []string {
    var result []string
    for _, str := range s {
        if str != "" {
            result = append(result, str)
        }
    }
    return result
}


// ContainsSubSlice returns the index of the first occurrence of the subSlice
// in the slice
// -1 is returned if the subSlice is not found
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
