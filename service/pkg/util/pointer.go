package util

func Int16Ptr(i int16) *int16 {
	return &i
}

func Uint16Ptr(i uint64) *uint64 {
	return &i
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
