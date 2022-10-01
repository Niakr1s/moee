package recorder

import (
	"bytes"
	"fmt"
	"time"
)

type Track struct {
	Extension string
	Start     time.Time
	End       time.Time

	Raw bytes.Buffer
}

func (t Track) String() string {
	return fmt.Sprintf("Track with ext %s, start: %v, end: %v, buffer len: %d", t.Extension, t.Start, t.End, t.Raw.Len())
}

func (t Track) HasData() bool {
	return t.Raw.Len() > 0
}
