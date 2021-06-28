package util

func Int16Ptr(i int16) *int16 {
	return &i
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
