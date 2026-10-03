package vsgo

import (
	"bufio"
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
)

var y4mSubsampling = map[[2]int]string{
	{1, 1}: "420", {1, 0}: "422", {0, 0}: "444",
	{2, 2}: "410", {2, 0}: "411", {0, 1}: "440",
}

func y4mHeader(vi VideoInfo) (string, error) {
	f := vi.Format
	format := y4mSubsampling[[2]int{f.SubSamplingW, f.SubSamplingH}]
	if f.ColorFamily != YUV || f.SampleType != integer || format == "" || vi.Width == 0 {
		return "", errors.New("vsgo: y4m needs a constant integer YUV clip")
	}
	if f.BitsPerSample > 8 {
		format += fmt.Sprintf("p%d", f.BitsPerSample)
	}
	return fmt.Sprintf("YUV4MPEG2 C%s W%d H%d F%d:%d Ip A0:0 XLENGTH=%d\n",
		format, vi.Width, vi.Height, vi.FPSNum, vi.FPSDen, vi.NumFrames), nil
}

// bufio keeps the first write error and flush returns it
func writeFrame(w *bufio.Writer, f *Frame, format VideoFormat) error {
	w.WriteString("FRAME\n")
	for plane := range format.NumPlanes {
		data, stride := f.Plane(plane), f.Stride(plane)
		rowSize := f.Width(plane) * format.BytesPerSample
		for y := range f.Height(plane) {
			w.Write(data[y*stride:][:rowSize])
		}
	}
	return w.Flush()
}

// same output as vspipe --container y4m, with `requests` frames in work at a time
func (n *Node) WriteY4M(ctx context.Context, w io.Writer, requests int) error {
	vi := n.Info()
	header, err := y4mHeader(vi)
	if err != nil {
		return err
	}
	bw := bufio.NewWriterSize(w, 1<<20) // 1 mib
	bw.WriteString(header)

	type result struct {
		frame *Frame
		err   error
	}
	requests = max(requests, 1)
	// frame i uses slot i%requests
	slots := make([]chan result, requests)
	for i := range slots {
		slots[i] = make(chan result, 1)
	}
	next := 0
	request := func() {
		if next == vi.NumFrames {
			return
		}
		i := next
		next++
		go func() {
			f, err := n.Frame(i)
			slots[i%requests] <- result{f, err}
		}()
	}
	for range requests {
		request()
	}

	for i := 0; i < next; i++ {
		res := <-slots[i%requests]
		err = cmp.Or(err, res.err, ctx.Err())
		if err == nil {
			err = writeFrame(bw, res.frame, vi.Format)
		}
		if res.frame != nil {
			res.frame.Free()
		}
		if err == nil {
			request()
		}
	}
	return cmp.Or(err, bw.Flush())
}
