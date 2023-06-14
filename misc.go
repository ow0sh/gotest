package main

func MapToSet(inter interface{}) map[string]struct{} {
	var result = make(map[string]struct{})
	switch val := inter.(type) {
	case []string:
		for _, str := range val {
			result[str] = struct{}{}
		}
	case map[string]string:
		for str := range val {
			result[str] = struct{}{}
		}
	default:
		panic("such type is not impl")
	}

	return result
}
