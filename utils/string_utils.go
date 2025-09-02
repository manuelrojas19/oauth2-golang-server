package utils

func StringDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
