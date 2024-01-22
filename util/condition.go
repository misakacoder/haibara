package util

func AllNotNull(objects ...any) bool {
	for _, object := range objects {
		str, ok := object.(string)
		if ok && str == "" {
			return false
		} else if object == nil {
			return false
		}
	}
	return true
}

func ConditionalExpression[T any](condition bool, value, defaultValue T) T {
	if condition {
		return value
	} else {
		return defaultValue
	}
}

func RequireNonNullElse[T any](value, defaultValue T) T {
	return ConditionalExpression(AllNotNull(value), value, defaultValue)
}
