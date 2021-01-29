package maputil

import "fmt"

// GetStringFromMap возвращает строковое значение по ключу
func GetStringFromMap(m map[string]interface{}, key string, defaultValue *string) (string, error) {
	v, ok := m[key]
	if !ok && defaultValue == nil {
		return "", fmt.Errorf("Не задано значение для поля %s", key)
	} else if !ok && defaultValue != nil {
		return *defaultValue, nil
	}

	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("Значение поля %s не является строкой", key)
	}

	return s, nil
}

// GetIntFromMap возвращает целочисленное значение по ключу
func GetIntFromMap(m map[string]interface{}, key string, defaultValue *int) (int, error) {
	v, ok := m[key]
	if !ok && defaultValue == nil {
		return 0, fmt.Errorf("Не задано значение для поля %s", key)
	} else if !ok && defaultValue != nil {
		return *defaultValue, nil
	}

	n, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("Значение поля %s не является числом", key)
	}

	if n != float64(int(n)) {
		return 0, fmt.Errorf("Значение поля %s не является целым числом", key)
	}

	return int(n), nil
}
