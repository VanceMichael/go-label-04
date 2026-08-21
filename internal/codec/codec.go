package codec

import (
	"compress/gzip"
	"encoding/json"
	"github.com/VanceMichael/go-base-airbridge/internal/domain"
	"io"
)

func Encode(v any) ([]byte, error) { return json.Marshal(v) }
func Decode(data []byte, v any) error {
	if len(data) == 0 || v == nil {
		return domain.ErrInvalid
	}
	return json.Unmarshal(data, v)
}
func Compress(data []byte) ([]byte, error) {
	var out []byte
	w := newSliceWriter(&out)
	gz := gzip.NewWriter(w)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return out, nil
}

type sliceWriter struct{ out *[]byte }

func newSliceWriter(out *[]byte) *sliceWriter      { return &sliceWriter{out: out} }
func (w *sliceWriter) Write(p []byte) (int, error) { *w.out = append(*w.out, p...); return len(p), nil }

var _ io.Writer = (*sliceWriter)(nil)
