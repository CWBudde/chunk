# Copilot Instructions for go-audio/chunk

This codebase provides a utility for reading structured data chunks, typically used in audio formats like WAV or AIFF.

## Architecture & Design

- **`Reader` Component**: The core of the library is [chunk.go](chunk.go). A `Reader` wraps an `io.Reader` and adds metadata like `ID` ([4]byte) and `Size`.
- **Position Tracking**: The `Reader` manually tracks progress via the `Pos` field. Always update or refer to `ch.Pos` when implementing new reading methods.
- **Resource Management**: The `Done()` method is critical. It must be called to ensure the underlying stream is drained to the end of the chunk, allowing subsequent chunks to be read from the same stream.

## Coding Patterns

- **Endianness**: Use `ReadLE` and `ReadBE` for struct/primitive reading. These methods wrap `binary.Read`.
- **Drain/Jump Pattern**: Use `io.CopyN(ioutil.Discard, ch.R, ...)` for skipping data.
- **Error Handling**: Follow the pattern of checking for nil `Reader` or nil inner `io.Reader` at the start of methods:
  ```go
  if ch == nil || ch.R == nil {
      return errors.New("nil Reader/reader pointer")
  }
  ```
- **EOF Handling**: `IsFullyRead()` should be checked before reading. Return `io.EOF` if the position has reached or exceeded the size.

## Testing Conventions

- Use `t.Run` for descriptive subtests in [chunk_test.go](chunk_test.go).
- Mock data using `bytes.NewReader`.
- Always verify that `r.Pos` is correctly incremented after a read or jump operation.

## Dependencies

- Standard library only (`io`, `encoding/binary`, etc.).
- Note: The codebase currently uses `io/ioutil` for `ioutil.Discard`. While deprecated in Go 1.16+, maintain consistency with existing code unless refactoring the entire project.
