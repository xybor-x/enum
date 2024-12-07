package mtmap

var globalmap = &MTMap{}

func Get[V any](key mtKeyer[V]) (V, bool) {
	return GetM(globalmap, key)
}

func MustGet[V any](key mtKeyer[V]) V {
	return MustGetM(globalmap, key)
}

func Set[V any](key mtKeyer[V], v V) {
	SetM(globalmap, key, v)
}
