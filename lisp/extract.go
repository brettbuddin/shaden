package lisp

// Extract performs a type assertion on a value and returns a user-friendly
// error on failure. argPos is 1-indexed for error messages.
func Extract[T any](v any, argPos int, typeName string) (T, error) {
	out, ok := v.(T)
	if !ok {
		var zero T
		return zero, ArgExpectError(typeName, argPos)
	}
	return out, nil
}

// ExtractArg extracts a typed value from args[i] with bounds checking.
// argPos in error messages is 1-indexed (i+1).
func ExtractArg[T any](args List, i int, typeName string) (T, error) {
	if i >= len(args) {
		var zero T
		return zero, ArgExpectError(typeName, i+1)
	}
	return Extract[T](args[i], i+1, typeName)
}

// ExtractFloat64 extracts a float64 from a value that may be int or float64.
func ExtractFloat64(v any, argPos int) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, ArgExpectError(AcceptTypes(TypeInt, TypeFloat), argPos)
	}
}

// ForEachEntry iterates over a List or Table, calling fn(key, value) for each
// entry. For Lists the key is the integer index; for Tables it is the map key.
// argPos is 1-indexed for the error message if the container type is wrong.
func ForEachEntry(container any, argPos int, fn func(key, value any) error) error {
	switch v := container.(type) {
	case List:
		for i, e := range v {
			if err := fn(i, e); err != nil {
				return err
			}
		}
	case Table:
		for k, e := range v {
			if err := fn(k, e); err != nil {
				return err
			}
		}
	default:
		return ArgExpectError(AcceptTypes(TypeList, TypeTable), argPos)
	}
	return nil
}
