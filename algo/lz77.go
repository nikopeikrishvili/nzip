package algo

import (
	"bytes"
	"encoding/gob"
	"os"
)

// LZ77Token represents a single token in the LZ77 encoded output
type LZ77Token struct {
	I bool // if its just character without match
	O int  // How far back to look
	L int  // Length of the match
	N byte // The next unmatched byte
}

func Compress(data []byte, windowSize int) []LZ77Token {
	var tokens []LZ77Token
	cursor := 0

	for cursor < len(data) {
		matchOffset, matchLength := 0, 0

		// Define the sliding window
		start := max(0, cursor-windowSize)

		// Search for the longest match in the sliding window
		for i := start; i < cursor; i++ {
			length := 0
			for cursor+length < len(data) && data[i+length] == data[cursor+length] {
				length++
				if i+length >= cursor {
					break
				}
			}

			if length > matchLength {
				matchOffset = cursor - i
				matchLength = length
			}
		}

		// Determine if this is a literal or a match
		if matchLength == 0 {
			// No match: store as a literal
			tokens = append(tokens, LZ77Token{
				I: true, // Mark as literal
				N: data[cursor],
			})
		} else {
			// Match found: store as a match
			next := byte(0)
			if cursor+matchLength < len(data) {
				next = data[cursor+matchLength]
			}
			tokens = append(tokens, LZ77Token{
				I: false,       // Mark as match
				O: matchOffset, // Offset
				L: matchLength, // Length
				N: next,        // Next unmatched byte
			})
		}

		// Move the cursor forward
		cursor += matchLength + 1
	}

	return tokens
}

// Decompress implements LZ77 decompression
func Decompress(compressedFilePath string, originalFilePath string) error {
	tokens, err := ReadCompressedFromFile(compressedFilePath)
	if err != nil {
		return err
	}
	// Write the decompressed data to the file
	err = os.WriteFile(originalFilePath, tokens, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Utility function to calculate the max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// WriteCompressedToFile writes compressed tokens to a file
func WriteCompressedToFile(tokens []LZ77Token, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	for _, token := range tokens {
		if token.I {
			// Write marker for literal

			// Write the literal byte
			if err := encoder.Encode(token.N); err != nil {
				return err
			}
		} else {

			//// Write the full token
			//if err := encoder.Encode(token); err != nil {
			//	return err
			//}
		}
	}
	return nil
}
func ReadCompressedFromFile(filePath string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a gob decoder and decode the tokens
	var tokens []LZ77Token
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&tokens); err != nil {
		return nil, err
	}

	// Reconstruct the original data
	var result bytes.Buffer
	for _, token := range tokens {
		if token.I {
			//	// If it's a literal, write the character directly
			result.WriteByte(token.N)
		} else {
			// If it's a match, reconstruct from the sliding window
			start := result.Len() - token.O
			for i := 0; i < token.L; i++ {
				result.WriteByte(result.Bytes()[start+i])
			}
			// Add the next unmatched byte
			result.WriteByte(token.N)
		}
	}

	return result.Bytes(), nil
}
