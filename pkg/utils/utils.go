package utils

func IntToPointerInt(v int64) *int64 {
	result := int64(v)
	return &result
}
func IntToPointerFloat(v uint64) *float64 {
	result := float64(v)
	return &result
}
func FloatToPointerFloat(v float64) *float64 {
	return &v
}
func FloatToPointerInt(v int64) *int64 {
	return &v
}
