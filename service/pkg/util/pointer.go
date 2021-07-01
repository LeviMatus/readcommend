package util

// Int16Ptr accepts an int16 and returns a pointer to that int16.
func Int16Ptr(i int16) *int16 {
	return &i
}

// Uint64Ptr accepts an uint16 and returns a pointer to that uint16.
func Uint64Ptr(i uint64) *uint64 {
	return &i
}

// StringPtr accepts a string and returns a pointer to that string.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
