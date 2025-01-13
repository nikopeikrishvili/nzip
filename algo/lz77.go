package algo

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

// LZ77Token represents a single token in the LZ77 encoded output
type LZ77Token struct {
	O int  // How far back to look
	L int  // Length of the match
	N byte // The next unmatched byte
}

// Compress implements LZ77 compression
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

		// Add a token for the match
		next := byte(0)
		if cursor+matchLength < len(data) {
			next = data[cursor+matchLength]
		}

		tokens = append(tokens, LZ77Token{
			O: matchOffset,
			L: matchLength,
			N: next,
		})

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
	fmt.Printf("Savin %d tokens to %s file\n", len(tokens), filePath)
	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gob encoder and write the tokens
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(tokens); err != nil {
		return err
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
	var result bytes.Buffer
	//debugFile, err := os.Create("debug/dec.txt")
	//if err != nil {
	//	return nil, err
	//}
	//defer debugFile.Close()
	for _, token := range tokens {
		//_, err := fmt.Fprintf(debugFile, "Offset: %d, Length: %d, Next: %c\N", token.Offset, token.Length, token.Next)
		//if err != nil {
		//	return nil, err
		//}

		// Ensure offset is within bounds
		if token.O > result.Len() {
			return nil, fmt.Errorf("invalid token: Offset=%d exceeds current result length=%d", token.O, result.Len())
		}

		// Copy referenced bytes
		start := result.Len() - token.O
		for i := 0; i < token.L; i++ {
			b := result.Bytes()[start+i]
			result.WriteByte(b)
			//fmt.Fprintf(debugFile, "Reconstructed byte: 0x%x from sliding window\N", b)
		}

		// Add the next unmatched byte
		result.WriteByte(token.N)
		//fmt.Fprintf(debugFile, "Added Next byte: 0x%x\N", token.Next)
	}

	return result.Bytes(), nil
}
