package mtmap

var globalmap = &MTMap{}

func Get2[V any](key mtKeyer[V]) (V, bool) {
	return Get2M(globalmap, key)
}

func Get[V any](key mtKeyer[V]) V {
	return GetM(globalmap, key)
}

func Set[V any](key mtKeyer[V], v V) {
	SetM(globalmap, key, v)
}
