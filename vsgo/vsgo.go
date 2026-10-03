package vsgo

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo linux LDFLAGS: -ldl
#include <stdlib.h>
#include "shim.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"
)

const (
	LogDebug = iota
	LogInformation
	LogWarning
	LogCritical
	LogFatal
)

var load = sync.OnceValue(func() error {
	names := []string{"libvapoursynth.so", "libvapoursynth.so.4", "libvapoursynth.so.0"}
	if runtime.GOOS == "windows" {
		names = []string{"vapoursynth.dll"}
	}
	var msg *C.char
	for _, name := range names {
		cname := C.CString(name)
		msg = C.vsgo_load(cname)
		C.free(unsafe.Pointer(cname))
		if msg == nil {
			return nil
		}
	}
	return fmt.Errorf("vsgo: cannot load VapourSynth: %s", C.GoString(msg))
})

type Core struct{ p *C.VSCore }

func NewCore() (*Core, error) {
	if err := load(); err != nil {
		return nil, err
	}
	return &Core{C.vsgo_createCore()}, nil
}

func (c *Core) Free() {
	C.vsgo_freeCore(c.p)
}

func (c *Core) Threads() int {
	var info C.VSCoreInfo
	C.vsgo_getCoreInfo(c.p, &info)
	return int(info.numThreads)
}

func (c *Core) OnLog(fn func(level int, msg string)) {
	C.vsgo_addLogHandler(c.p, C.uintptr_t(cgo.NewHandle(fn)))
}

//export vsgoLog
func vsgoLog(level C.int, msg *C.char, handle C.uintptr_t) {
	cgo.Handle(handle).Value().(func(int, string))(int(level), C.GoString(msg))
}

//export vsgoLogFree
func vsgoLogFree(handle C.uintptr_t) {
	cgo.Handle(handle).Delete()
}

func (c *Core) Invoke(namespace, function string, args Args) (*Map, error) {
	cns, cfn := C.CString(namespace), C.CString(function)
	defer C.free(unsafe.Pointer(cns))
	defer C.free(unsafe.Pointer(cfn))

	plugin := C.vsgo_getPluginByNamespace(cns, c.p)
	if plugin == nil {
		return nil, fmt.Errorf("vsgo: plugin %q is not loaded", namespace)
	}
	in, err := newArgs(args)
	if err != nil {
		return nil, fmt.Errorf("vsgo: %s.%s: %w", namespace, function, err)
	}
	defer C.vsgo_freeMap(in)

	out := C.vsgo_invoke(plugin, cfn, in)
	if msg := C.vsgo_mapGetError(out); msg != nil {
		defer C.vsgo_freeMap(out)
		return nil, fmt.Errorf("vsgo: %s.%s: %s", namespace, function, C.GoString(msg))
	}
	return &Map{out}, nil
}

func (c *Core) Clip(namespace, function string, args Args) (*Node, error) {
	m, err := c.Invoke(namespace, function, args)
	if err != nil {
		return nil, err
	}
	defer m.Free()
	n := m.Node("clip")
	if n == nil {
		return nil, fmt.Errorf("vsgo: %s.%s returned no clip", namespace, function)
	}
	return n, nil
}

func (c *Core) LoadPlugin(path string) error {
	m, err := c.Invoke("std", "LoadPlugin", Args{"path": path})
	if err != nil {
		return err
	}
	m.Free()
	return nil
}
