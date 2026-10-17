package vsgo

/*
#include <stdlib.h>
#include "shim.h"
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

type Args map[string]any

func newArgs(args Args) (*C.VSMap, error) {
	m := C.vsgo_createMap()
	for key, value := range args {
		ckey := C.CString(key)
		err := setArg(m, ckey, value)
		C.free(unsafe.Pointer(ckey))
		if err != nil {
			C.vsgo_freeMap(m)
			return nil, fmt.Errorf("argument %q: %w", key, err)
		}
	}
	return m, nil
}

func setArg(m *C.VSMap, key *C.char, value any) error {
	switch v := value.(type) {
	case bool:
		i := 0
		if v {
			i = 1
		}
		C.vsgo_mapSetInt(m, key, C.int64_t(i))
	case int:
		C.vsgo_mapSetInt(m, key, C.int64_t(v))
	case int64:
		C.vsgo_mapSetInt(m, key, C.int64_t(v))
	case float64:
		C.vsgo_mapSetFloat(m, key, C.double(v))
	case string:
		C.vsgo_mapSetString(m, key, (*C.char)(unsafe.Pointer(unsafe.StringData(v))), C.int(len(v)))
	case *Node:
		C.vsgo_mapSetNode(m, key, v.p)
		runtime.KeepAlive(v)
	default:
		slice := reflect.ValueOf(value)
		if slice.Kind() != reflect.Slice {
			return fmt.Errorf("unsupported type %T", value)
		}
		for i := range slice.Len() {
			if err := setArg(m, key, slice.Index(i).Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func mapInt(m *C.VSMap, key string) (int64, bool) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var e C.int
	v := C.vsgo_mapGetInt(m, ckey, &e)
	return int64(v), e == 0
}

type Map struct{ p *C.VSMap }

func (m *Map) Free() {
	C.vsgo_freeMap(m.p)
}

func (m *Map) Int(key string) (int64, bool) {
	return mapInt(m.p, key)
}

func (m *Map) Node(key string) *Node {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	p := C.vsgo_mapGetNode(m.p, ckey)
	if p == nil {
		return nil
	}
	n := &Node{p}
	runtime.AddCleanup(n, func(p *C.VSNode) { C.vsgo_freeNode(p) }, p)
	return n
}
