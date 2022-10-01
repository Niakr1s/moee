package mp3

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"strconv"
)

var (
	ErrNoMetaint         = errors.New("mp3: key 'icy-metaint' not found in HTTP header")
	ErrCorruptedMetadata = errors.New("mp3: corrupted metadata")
	ErrNoStreamTitle     = errors.New("mp3: no 'StreamTitle' tag in metadata")
)

type Extractor struct {
	metaint int64 // Distance between two metadata chunks
}

func NewExtractor(respHdr http.Header) (*Extractor, error) {
	mi := respHdr.Get("icy-metaint")
	if mi == "" {
		return nil, ErrNoMetaint
	}
	miNum, _ := strconv.ParseInt(mi, 10, 64)
	return &Extractor{
		metaint: miNum,
	}, nil
}

func (d *Extractor) ReadBlock(r io.Reader, w io.Writer) (isFirst bool, err error) {
	var musicData bytes.Buffer

	// We want to write everything except the metadata to the output and to
	// musicData for calculating the checksum.
	multi := io.MultiWriter(w, &musicData)

	// Read until the metadata chunk. The part that is read here is also what
	// contains the actual mp3 music data.
	io.CopyN(multi, r, d.metaint)

	// Read number of metadata blocks (blocks within this function are not what
	// is meant with `ReadBlock()`).
	var numBlocks uint8
	err = binary.Read(r, binary.LittleEndian, &numBlocks)
	if err != nil {
		return false, err
	}

	return numBlocks > 0, nil
}

func (d *Extractor) GetExtension() string {
	return ".mp3"
}
