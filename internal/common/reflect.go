package common

import "reflect"

func NameOf[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().Name()
}
