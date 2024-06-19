package main

import "fmt"

type Flags struct {
	m map[string]any
}

type FlagValue interface {
	bool | int64 | string
}

// type Value[T FlagValue] struct {
// 	value *T
// }
// type Flag[V FlagValue] struct {
// 	Name  string
// 	Value V
// }

// func GetValue[V FlagValue](name string, flags Flag[V]) {
// 	// var m map[string]FlagValue doesn't work
// 	// var m map[string]interface{}
// }

func SetFlag[V FlagValue](flags *Flags, name string, v V) {
	flags.m[name] = v
}

func GetFlag[V FlagValue](flags *Flags, name string) (V, bool) {
	v := flags.m[name]
	vv, ok := v.(V)
	return vv, ok
}

func main() {
	flagsMap := make(map[string]any)
	flags := Flags{flagsMap}
	SetFlag(&flags, "foo", int64(42))
	SetFlag(&flags,	"bar", "deadbeef")
	foo, ok := GetFlag[int64](&flags, "foo")
	if !ok {
		fmt.Println("not ok")
	} else {
		fmt.Printf("foo is type %T value %v\n", foo, foo)
	}
	bar, ok := GetFlag[string](&flags, "bar")
	if !ok {
		fmt.Println("not ok")
	} else {
		fmt.Printf("bar is type %T value %v\n", bar, bar)
	}
}
