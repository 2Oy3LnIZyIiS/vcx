package stringkit

import "strings"


func Equals(a, b string, ignoreCase bool) bool {
    if ignoreCase {
        return strings.EqualFold(a, b)
    }
    return a == b
}


// containsComponent checks if one component contains another, respecting options
func Contains(a, b string, ignoreCase bool) bool {
    if ignoreCase {
        return strings.Contains(
            strings.ToLower(a),
            strings.ToLower(b),
        )
    }
    return strings.Contains(a, b)
}
