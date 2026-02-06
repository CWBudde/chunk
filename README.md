# Chunk

A Go library for reading structured data chunks found in audio file formats like WAV (RIFF) and AIFF.

## Install

```bash
go get github.com/CWBudde/chunk
```

## Usage

The `Reader` wraps an `io.Reader` with chunk metadata -- a 4-byte ID, size, and position tracking. It's designed to be used by a parent container parser that reads chunk headers and hands off the body to a `Reader`.

```go
// Typically created by a container parser (RIFF, AIFF, etc.)
ch := &chunk.Reader{
    ID:   [4]byte{'f', 'm', 't', ' '},
    Size: 16,
    R:    underlyingReader,
}

// Read typed data in little or big endian
var sampleRate uint32
ch.ReadLE(&sampleRate)

// Read a single byte
b, err := ch.ReadByte()

// Skip ahead
ch.Jump(4)

// Check if all data has been consumed
if ch.IsFullyRead() {
    // ...
}

// IMPORTANT: always call Done() when finished to drain
// any unread bytes and advance the stream for the next chunk
ch.Done()
```

## API

| Method | Description |
| --- | --- |
| `Read(p []byte)` | Implements `io.Reader` |
| `ReadLE(dst any)` | Read into `dst` using little-endian byte order |
| `ReadBE(dst any)` | Read into `dst` using big-endian byte order |
| `ReadByte()` | Read and return a single byte |
| `Jump(n int)` | Skip ahead `n` bytes |
| `IsFullyRead()` | Returns true if position >= size |
| `Done()` | Drains any remaining unread bytes |

## License

Apache 2.0 -- see [LICENSE](LICENSE).
