package inmemorydb

func compareNumeric(value interface{}, conditionValue interface{}, comparator func(float64, float64) bool) bool {
	valueFloat, ok1 := convertToFloat(value)
	conditionFloat, ok2 := convertToFloat(conditionValue)
	if !ok1 || !ok2 {
		return false
	}
	return comparator(valueFloat, conditionFloat)
}

func convertToFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	default:
		return 0, false
	}
}
