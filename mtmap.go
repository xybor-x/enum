package enum

type mtmap struct {
	data map[any]any
}

type keyI[T any] interface {
	getValue() T
}

func get2[V any](m *mtmap, key keyI[V]) (V, bool) {
	var zero V
	if m.data == nil {
		return zero, false
	}
	val, exists := m.data[key]
	if !exists {
		return zero, false
	}

	// Type assertion can only fail if V is an interface. In that
	// case, if the map has a `nil` in it, Go won't be able to
	// type assert that nil into the interface value. (Note a nil
	// *pointer* type asserts just fine, because it is still
	// carrying a concrete type. Only nil interfaces lack any
	// concrete type.) So if the type assertion fails, it must be
	// a nil, here played by the zero we declared above, which
	// must be `nil` even though the compiler can't realize that
	// `nil` would be safe here.
	finalVal, canAssert := val.(V)
	if canAssert {
		return finalVal, true
	} else {
		return zero, true
	}
}

func get[V any](m *mtmap, key keyI[V]) V {
	v, _ := get2(m, key)
	return v
}

func set[V any](m *mtmap, key keyI[V], val V) {
	if m.data == nil {
		m.data = map[any]any{}
	}

	m.data[key] = val
}

type key[Key any, Value any] struct {
	key Key
}

// this is not called, but it implements the type restriction
func (k key[K, V]) getValue() (val V) { return }
