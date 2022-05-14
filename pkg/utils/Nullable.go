package utils

type Nullable[T interface{}] struct {
	value    T
	hasValue bool
}

func (val *Nullable[T]) SetValue(v T) {
	val.value = v
	val.hasValue = true
}

func (val *Nullable[T]) GetValue() T {
	if val.hasValue {
		return val.value
	}

	panic("No value")
}

func (val *Nullable[T]) HasValue() bool {
	return val.hasValue
}
