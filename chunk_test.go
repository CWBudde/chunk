package chunk

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func TestReader_Read(t *testing.T) {
	t.Run("reads data and advances position", func(t *testing.T) {
		data := []byte("hello")
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		buf := make([]byte, 3)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 3 {
			t.Fatalf("expected n=3, got %d", n)
		}
		if r.Pos != 3 {
			t.Fatalf("expected Pos=3, got %d", r.Pos)
		}
		if string(buf) != "hel" {
			t.Fatalf("expected 'hel', got %q", buf)
		}
	})

	t.Run("reads remaining data and returns EOF", func(t *testing.T) {
		data := []byte("hi")
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		buf := make([]byte, 5)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error on first read: %v", err)
		}
		if n != 2 {
			t.Fatalf("expected n=2, got %d", n)
		}

		n, err = r.Read(buf)
		if err != io.EOF {
			t.Fatalf("expected EOF, got %v", err)
		}
		if n != 0 {
			t.Fatalf("expected n=0 on EOF, got %d", n)
		}
	})

	t.Run("nil reader returns error", func(t *testing.T) {
		r := &Reader{}
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err == nil {
			t.Fatal("expected error for nil reader")
		}
	})

	t.Run("nil Reader pointer returns error", func(t *testing.T) {
		var r *Reader
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err == nil {
			t.Fatal("expected error for nil Reader pointer")
		}
	})
}

func TestReader_ReadLE(t *testing.T) {
	t.Run("reads uint16 little endian", func(t *testing.T) {
		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, 0x0102)
		r := &Reader{
			Size: 2,
			R:    bytes.NewReader(buf),
		}

		var val uint16
		err := r.ReadLE(&val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 0x0102 {
			t.Fatalf("expected 0x0102, got 0x%04x", val)
		}
		if r.Pos != 2 {
			t.Fatalf("expected Pos=2, got %d", r.Pos)
		}
	})

	t.Run("reads uint32 little endian", func(t *testing.T) {
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, 42)
		r := &Reader{
			Size: 4,
			R:    bytes.NewReader(buf),
		}

		var val uint32
		err := r.ReadLE(&val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 42 {
			t.Fatalf("expected 42, got %d", val)
		}
	})
}

func TestReader_ReadBE(t *testing.T) {
	t.Run("reads uint16 big endian", func(t *testing.T) {
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, 0x0102)
		r := &Reader{
			Size: 2,
			R:    bytes.NewReader(buf),
		}

		var val uint16
		err := r.ReadBE(&val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 0x0102 {
			t.Fatalf("expected 0x0102, got 0x%04x", val)
		}
	})

	t.Run("reads uint32 big endian", func(t *testing.T) {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, 12345)
		r := &Reader{
			Size: 4,
			R:    bytes.NewReader(buf),
		}

		var val uint32
		err := r.ReadBE(&val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 12345 {
			t.Fatalf("expected 12345, got %d", val)
		}
	})
}

func TestReader_ReadByte(t *testing.T) {
	t.Run("reads single byte", func(t *testing.T) {
		r := &Reader{
			Size: 3,
			R:    bytes.NewReader([]byte{0xAA, 0xBB, 0xCC}),
		}

		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b != 0xAA {
			t.Fatalf("expected 0xAA, got 0x%02x", b)
		}
	})

	t.Run("reads multiple bytes sequentially", func(t *testing.T) {
		r := &Reader{
			Size: 2,
			R:    bytes.NewReader([]byte{0x01, 0x02}),
		}

		b1, err := r.ReadByte()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		b2, err := r.ReadByte()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b1 != 0x01 || b2 != 0x02 {
			t.Fatalf("expected 0x01, 0x02 got 0x%02x, 0x%02x", b1, b2)
		}
	})

	t.Run("returns EOF when fully read", func(t *testing.T) {
		r := &Reader{
			Size: 0,
			R:    bytes.NewReader(nil),
		}

		_, err := r.ReadByte()
		if err != io.EOF {
			t.Fatalf("expected EOF, got %v", err)
		}
	})

	t.Run("returns EOF after reading all bytes", func(t *testing.T) {
		r := &Reader{
			Size: 1,
			R:    bytes.NewReader([]byte{0xFF}),
		}

		_, err := r.ReadByte()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = r.ReadByte()
		if err != io.EOF {
			t.Fatalf("expected EOF, got %v", err)
		}
	})
}

func TestReader_IsFullyRead(t *testing.T) {
	t.Run("false when data remains", func(t *testing.T) {
		r := &Reader{
			Size: 10,
			R:    bytes.NewReader(make([]byte, 10)),
			Pos:  0,
		}
		if r.IsFullyRead() {
			t.Fatal("expected false, got true")
		}
	})

	t.Run("true when fully consumed", func(t *testing.T) {
		r := &Reader{
			Size: 5,
			R:    bytes.NewReader(make([]byte, 5)),
			Pos:  5,
		}
		if !r.IsFullyRead() {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("true when Pos exceeds Size", func(t *testing.T) {
		r := &Reader{
			Size: 3,
			R:    bytes.NewReader(make([]byte, 3)),
			Pos:  10,
		}
		if !r.IsFullyRead() {
			t.Fatal("expected true, got false")
		}
	})

	t.Run("true for nil Reader pointer", func(t *testing.T) {
		var r *Reader
		if !r.IsFullyRead() {
			t.Fatal("expected true for nil pointer")
		}
	})

	t.Run("true for nil inner reader", func(t *testing.T) {
		r := &Reader{Size: 5}
		if !r.IsFullyRead() {
			t.Fatal("expected true for nil inner reader")
		}
	})
}

func TestReader_Jump(t *testing.T) {
	t.Run("jumps ahead by N bytes", func(t *testing.T) {
		data := []byte("abcdefghij")
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		err := r.Jump(5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Pos != 5 {
			t.Fatalf("expected Pos=5, got %d", r.Pos)
		}

		// Read the next byte to verify position
		buf := make([]byte, 1)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 1 || buf[0] != 'f' {
			t.Fatalf("expected 'f', got %q", buf[0])
		}
	})

	t.Run("jump zero does nothing", func(t *testing.T) {
		r := &Reader{
			Size: 5,
			R:    bytes.NewReader([]byte("hello")),
		}

		err := r.Jump(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Pos != 0 {
			t.Fatalf("expected Pos=0, got %d", r.Pos)
		}
	})

	t.Run("jump negative does nothing", func(t *testing.T) {
		r := &Reader{
			Size: 5,
			R:    bytes.NewReader([]byte("hello")),
		}

		err := r.Jump(-3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Pos != 0 {
			t.Fatalf("expected Pos=0, got %d", r.Pos)
		}
	})

	t.Run("jump beyond data returns error", func(t *testing.T) {
		r := &Reader{
			Size: 3,
			R:    bytes.NewReader([]byte("abc")),
		}

		err := r.Jump(100)
		if err == nil {
			t.Fatal("expected error jumping beyond data")
		}
	})
}

func TestReader_Done(t *testing.T) {
	t.Run("drains remaining data", func(t *testing.T) {
		data := []byte("hello world")
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		// Read partial
		buf := make([]byte, 5)
		r.Read(buf)

		if err := r.Done(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns nil when already fully read", func(t *testing.T) {
		data := []byte("hi")
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
			Pos:  len(data),
		}

		if err := r.Done(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns nil for zero-size reader", func(t *testing.T) {
		r := &Reader{
			Size: 0,
			R:    bytes.NewReader(nil),
		}
		if err := r.Done(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestReader_readWithByteOrder(t *testing.T) {
	t.Run("returns error for nil reader pointer", func(t *testing.T) {
		var r *Reader
		var val uint16
		err := r.readWithByteOrder(&val, binary.LittleEndian)
		if err == nil {
			t.Fatal("expected error for nil Reader")
		}
	})

	t.Run("returns error for nil inner reader", func(t *testing.T) {
		r := &Reader{Size: 2}
		var val uint16
		err := r.readWithByteOrder(&val, binary.LittleEndian)
		if err == nil {
			t.Fatal("expected error for nil inner reader")
		}
	})

	t.Run("returns EOF when fully read", func(t *testing.T) {
		r := &Reader{
			Size: 2,
			R:    bytes.NewReader([]byte{0, 0}),
			Pos:  2,
		}
		var val uint16
		err := r.readWithByteOrder(&val, binary.LittleEndian)
		if err != io.EOF {
			t.Fatalf("expected EOF, got %v", err)
		}
	})

	t.Run("returns error on short underlying reader", func(t *testing.T) {
		// Size says 4 bytes available, but underlying reader only has 1 byte
		r := &Reader{
			Size: 4,
			R:    bytes.NewReader([]byte{0x01}),
		}
		var val uint32
		err := r.readWithByteOrder(&val, binary.LittleEndian)
		if err == nil {
			t.Fatal("expected error for short reader")
		}
		// Pos should NOT have advanced since binary.Read failed
		if r.Pos != 0 {
			t.Fatalf("expected Pos=0 after failed read, got %d", r.Pos)
		}
	})
}

func TestReader_drain(t *testing.T) {
	t.Run("drains remaining bytes", func(t *testing.T) {
		data := []byte("abcdef")
		underlying := bytes.NewReader(data)
		r := &Reader{
			Size: len(data),
			R:    underlying,
			Pos:  2,
		}

		err := r.drain()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("noop when nothing to drain", func(t *testing.T) {
		r := &Reader{
			Size: 5,
			R:    bytes.NewReader([]byte("hello")),
			Pos:  5,
		}

		err := r.drain()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("noop when Pos exceeds Size", func(t *testing.T) {
		r := &Reader{
			Size: 3,
			R:    bytes.NewReader([]byte("abc")),
			Pos:  10,
		}

		err := r.drain()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestReader_ID(t *testing.T) {
	r := &Reader{
		ID:   [4]byte{'R', 'I', 'F', 'F'},
		Size: 0,
		R:    bytes.NewReader(nil),
	}
	if r.ID != [4]byte{'R', 'I', 'F', 'F'} {
		t.Fatalf("expected RIFF, got %s", r.ID[:])
	}
}

func TestReader_Integration(t *testing.T) {
	t.Run("read mixed types sequentially", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, uint16(1000))
		binary.Write(&buf, binary.LittleEndian, uint32(2000))
		binary.Write(&buf, binary.LittleEndian, byte(0xFF))

		data := buf.Bytes()
		r := &Reader{
			ID:   [4]byte{'d', 'a', 't', 'a'},
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		var u16 uint16
		if err := r.ReadLE(&u16); err != nil {
			t.Fatalf("ReadLE uint16: %v", err)
		}
		if u16 != 1000 {
			t.Fatalf("expected 1000, got %d", u16)
		}

		var u32 uint32
		if err := r.ReadLE(&u32); err != nil {
			t.Fatalf("ReadLE uint32: %v", err)
		}
		if u32 != 2000 {
			t.Fatalf("expected 2000, got %d", u32)
		}

		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("ReadByte: %v", err)
		}
		if b != 0xFF {
			t.Fatalf("expected 0xFF, got 0x%02x", b)
		}

		if !r.IsFullyRead() {
			t.Fatal("expected fully read after consuming all data")
		}
	})

	t.Run("jump then read", func(t *testing.T) {
		data := []byte{0x00, 0x00, 0x00, 0x00, 0xAB}
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		if err := r.Jump(4); err != nil {
			t.Fatalf("Jump: %v", err)
		}

		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("ReadByte after Jump: %v", err)
		}
		if b != 0xAB {
			t.Fatalf("expected 0xAB, got 0x%02x", b)
		}
	})

	t.Run("partial read then Done", func(t *testing.T) {
		data := make([]byte, 100)
		for i := range data {
			data[i] = byte(i)
		}
		r := &Reader{
			Size: len(data),
			R:    bytes.NewReader(data),
		}

		buf := make([]byte, 10)
		r.Read(buf)
		if r.Pos != 10 {
			t.Fatalf("expected Pos=10, got %d", r.Pos)
		}

		if err := r.Done(); err != nil {
			t.Fatalf("Done: %v", err)
		}
	})
}
