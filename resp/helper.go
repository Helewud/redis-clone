package resp

// Checks if a line ends with CRLF
func isEndOfLine(line []byte) bool {
	return len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n'
}

// Removes the trailing CRLF sequence
func trimEndOfLine(line []byte) []byte {
	if isEndOfLine(line) {
		return line[:len(line)-2]
	}
	return line
}
