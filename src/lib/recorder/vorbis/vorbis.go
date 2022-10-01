package vorbis

import (
	"bytes"
	"errors"
	"io"
)

var (
	ErrNoHeaderSegment = errors.New("vorbis: no header segment")
)

type Extractor struct {
	checksum uint32 // Used for an alternate filename when there's no metadata.
}

func NewExtractor() (*Extractor, error) {
	return new(Extractor), nil
}

func (d *Extractor) ReadBlock(reader io.Reader, w io.Writer) (isFirst bool, err error) {
	// Everything we read here is part of the music data so we can just use a
	// tee reader.
	r := io.TeeReader(reader, w)

	// Decode page.
	page, err := OggDecode(r)
	if err != nil {
		return false, err
	}

	// We need to be able to access `page.Segments[0]`.
	if len(page.Segments) == 0 {
		return false, ErrNoHeaderSegment
	}

	// Decode Vorbis header, stored in `page.Segments[0]`.
	hdr, err := VorbisHeaderDecode(bytes.NewBuffer(page.Segments[0]))
	if err != nil {
		return false, err
	}

	// Extract potential metadata.
	if hdr.PackType == PackTypeComment {
		d.checksum = page.Header.Checksum
	}

	// Return true for isFirst if this block is the beginning of a new file.
	return (page.Header.HeaderType & FHeaderTypeBOS) > 0, nil
}

func (d *Extractor) GetExtension() string {
	return ".ogg"
}
