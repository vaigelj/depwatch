package vulncheck

import "bytes"

// bytesReader wraps a byte slice in an io.Reader.
func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
