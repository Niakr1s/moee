package recorder

import (
	"io"
)

type Extractor interface {
	// Reads a single "block" from a radio stream. A block can be any chunk of
	// data, depending on the file format, for example in Ogg/Vorbis it would
	// be equivalent to a chunk. Writes the part containing the actual music
	// data into `w`.
	// `isFirst` is true, if the block read was the first block of a file.
	ReadBlock(r io.Reader, w io.Writer) (isFirst bool, err error)
	// for example, ".mp3" or ".ogg"
	GetExtension() string
}
