package main

// IsPathSeparator reports whether c is a directory separator character.
func IsPathSeparator(c uint8) bool {
	// NOTE: Windows accept / as path separator.
	return c == '\\' || c == '/'
}

func FileEnd(path string) int {
	for i := len(path) - 1; i >= 0 && !IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			part := path[i:]
			for ii := 0; ii < len(part); ii++ {
				if part[ii] == ' ' {
					return i + ii
				}
			}
		}
	}
	return -1
}
