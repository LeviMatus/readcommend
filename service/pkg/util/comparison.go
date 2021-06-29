package util

// Int16InRange compares the value of an int16 against the min
// and max values provided. If i is between these two, then
// true is returned.
func Int16InRange(i int16, min, max int16) bool {
	return min <= i && i <= max
}
