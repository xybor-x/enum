package mtmap

type MTMap struct {
	data map[any]any
}

type mtKeyer[V any] interface {
	InferValue() V
}

func Get2M[V any](m *MTMap, key mtKeyer[V]) (V, bool) {
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

func GetM[V any](m *MTMap, key mtKeyer[V]) V {
	v, _ := Get2M(m, key)
	return v
}

func SetM[V any](m *MTMap, key mtKeyer[V], val V) {
	if m.data == nil {
		m.data = map[any]any{}
	}

	m.data[key] = val
}
