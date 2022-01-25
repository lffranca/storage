package storage

import (
	"io"
	"time"
)

type File struct {
	Name         *string    `json:"name,omitempty" yaml:"name,omitempty"`
	Extension    *string    `json:"extension,omitempty" yaml:"extension,omitempty"`
	Key          *string    `json:"key,omitempty" yaml:"key,omitempty"`
	Path         *string    `json:"path,omitempty" yaml:"path,omitempty"`
	Size         *int       `json:"size,omitempty" yaml:"size,omitempty"`
	ContentType  *string    `json:"content_type,omitempty" yaml:"content_type,omitempty"`
	LastModified *time.Time `json:"last_modified,omitempty" yaml:"last_modified,omitempty"`
	Prefix       *string    `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

func NewBufferWriterAt(w io.Writer) *BufferWriterAt {
	buff := new(BufferWriterAt)
	buff.w = w

	return buff
}

type BufferWriterAt struct {
	w io.Writer
}

func (fw BufferWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}
