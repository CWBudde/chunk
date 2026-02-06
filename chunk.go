package chunk

import (
	"encoding/binary"
	"errors"
	"io"
)

// Reader is a struct representing a data chunk. Its reader is shared with the
// container but convenience methods are provided.
type Reader struct {
	ID   [4]byte
	Size int
	R    io.Reader
	Pos  int
}

// Done makes sure the entire Reader was read.
func (ch *Reader) Done() error {
	if !ch.IsFullyRead() {
		return ch.drain()
	}
	return nil
}

// Read implements the io.Reader interface.
func (ch *Reader) Read(p []byte) (n int, err error) {
	if ch == nil || ch.R == nil {
		return 0, errors.New("nil Reader/reader pointer")
	}
	n, err = ch.R.Read(p)
	ch.Pos += n
	return n, err
}

// ReadLE reads the Little Endian Reader data into the passed struct
func (ch *Reader) ReadLE(dst any) error {
	return ch.readWithByteOrder(dst, binary.LittleEndian)
}

// ReadBE reads the Big Endian Reader data into the passed struct
func (ch *Reader) ReadBE(dst any) error {
	return ch.readWithByteOrder(dst, binary.BigEndian)
}

// ReadByte reads and returns a single byte
func (ch *Reader) ReadByte() (byte, error) {
	if ch.IsFullyRead() {
		return 0, io.EOF
	}
	var b byte
	err := ch.ReadLE(&b)
	return b, err
}

// IsFullyRead checks if we're finished reading the Reader
func (ch *Reader) IsFullyRead() bool {
	if ch == nil || ch.R == nil {
		return true
	}
	return ch.Size <= ch.Pos
}

// Jump jumps ahead in the Reader
func (ch *Reader) Jump(bytesAhead int) error {
	var err error
	var n int64
	if bytesAhead > 0 {
		n, err = io.CopyN(io.Discard, ch.R, int64(bytesAhead))
		ch.Pos += int(n)
	}
	return err
}

func (ch *Reader) readWithByteOrder(dst any, byteOrder binary.ByteOrder) error {
	if ch == nil || ch.R == nil {
		return errors.New("nil Reader/reader pointer")
	}
	if ch.IsFullyRead() {
		return io.EOF
	}
	if err := binary.Read(ch.R, byteOrder, dst); err != nil {
		return err
	}
	ch.Pos += binary.Size(dst)
	return nil
}

// You are probably looking to call Done() instead!
func (ch *Reader) drain() error {
	bytesAhead := ch.Size - ch.Pos
	if bytesAhead > 0 {
		_, err := io.CopyN(io.Discard, ch.R, int64(bytesAhead))
		return err
	}
	return nil
}
