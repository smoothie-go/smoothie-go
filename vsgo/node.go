package vsgo

/*
#include "shim.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

const (
	Gray = C.cfGray
	RGB  = C.cfRGB
	YUV  = C.cfYUV

	integer = C.stInteger

	YUV444P8 = C.pfYUV444P8
	RGB48    = C.pfRGB48
)

type VideoFormat struct {
	ColorFamily    int
	SampleType     int
	BitsPerSample  int
	BytesPerSample int
	SubSamplingW   int
	SubSamplingH   int
	NumPlanes      int
}

func (f VideoFormat) ID() int {
	// same layout as VS_MAKE_VIDEO_ID
	return f.ColorFamily<<28 | f.SampleType<<24 | f.BitsPerSample<<16 | f.SubSamplingW<<8 | f.SubSamplingH
}

type VideoInfo struct {
	Format    VideoFormat
	FPSNum    int64
	FPSDen    int64
	Width     int
	Height    int
	NumFrames int
}

// freed by the garbage collector
type Node struct{ p *C.VSNode }

func (n *Node) Info() VideoInfo {
	vi := C.vsgo_getVideoInfo(n.p)
	defer runtime.KeepAlive(n)
	return VideoInfo{
		Format: VideoFormat{
			ColorFamily:    int(vi.format.colorFamily),
			SampleType:     int(vi.format.sampleType),
			BitsPerSample:  int(vi.format.bitsPerSample),
			BytesPerSample: int(vi.format.bytesPerSample),
			SubSamplingW:   int(vi.format.subSamplingW),
			SubSamplingH:   int(vi.format.subSamplingH),
			NumPlanes:      int(vi.format.numPlanes),
		},
		FPSNum:    int64(vi.fpsNum),
		FPSDen:    int64(vi.fpsDen),
		Width:     int(vi.width),
		Height:    int(vi.height),
		NumFrames: int(vi.numFrames),
	}
}

func (n *Node) Frame(i int) (*Frame, error) {
	var msg [1024]C.char
	p := C.vsgo_getFrame(C.int(i), n.p, &msg[0], C.int(len(msg)))
	runtime.KeepAlive(n)
	if p == nil {
		return nil, fmt.Errorf("vsgo: frame %d: %s", i, C.GoString(&msg[0]))
	}
	return &Frame{p}, nil
}

type Frame struct{ p *C.VSFrame }

func (f *Frame) Free() {
	C.vsgo_freeFrame(f.p)
}

func (f *Frame) Width(plane int) int {
	return int(C.vsgo_getFrameWidth(f.p, C.int(plane)))
}

func (f *Frame) Height(plane int) int {
	return int(C.vsgo_getFrameHeight(f.p, C.int(plane)))
}

func (f *Frame) Stride(plane int) int {
	return int(C.vsgo_getStride(f.p, C.int(plane)))
}

func (f *Frame) Plane(plane int) []byte {
	p := unsafe.Pointer(C.vsgo_getReadPtr(f.p, C.int(plane)))
	return unsafe.Slice((*byte)(p), f.Stride(plane)*f.Height(plane))
}

func (f *Frame) PropInt(key string) (int64, bool) {
	return mapInt(C.vsgo_getFrameProperties(f.p), key)
}
