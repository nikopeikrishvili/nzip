package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/nikopeikrishvili/nzip/algo"
	"io"
	"os"
	"path/filepath"
)

// CompressFile reads the input file, compresses its contents using RLE, and returns the compressed data
func CompressFile(filePath string, outputFilePath string) error {
	// Open the file in binary mode
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file contents into a buffer
	var inputBuffer bytes.Buffer
	_, err = io.Copy(&inputBuffer, file)
	if err != nil {
		return err
	}
	windowSize := 6
	compressed := algo.Compress(inputBuffer.Bytes(), windowSize)
	//debugFile, err := os.Create("debug/comp.txt")
	//if err != nil {
	//	return err
	//}
	//defer debugFile.Close()
	//for _, token := range compressed {
	//	_, err := fmt.Fprintf(debugFile, "Offset: %d, Length: %d, Next: %c\n", token.Offset, token.Length, token.Next)
	//	if err != nil {
	//		return err
	//	}
	//}
	fmt.Printf("Size of compressed tokens array %d\n", len(compressed))
	_ = algo.WriteCompressedToFile(compressed, outputFilePath)

	return nil
}

func DecompressFile(filePath string, outputFilePath string) error {

	err := algo.Decompress(filePath, outputFilePath)
	if err != nil {
		return err
	}
	return nil
}
func main() {
	// Define flags for file input and output
	mode := flag.String("mode", "compress", "compress or decompress")
	inputFile := flag.String("input", "", "Path to the input file to be compressed")
	outputFile := flag.String("output", "", "Path to the output file to save compressed data (optional)")

	flag.Parse()

	// Validate input
	if *inputFile == "" {
		fmt.Println("Error: input file path is required")
		flag.Usage()
		os.Exit(1)
	}
	// Determine the output file name
	if *outputFile == "" {
		// Default output file name
		*outputFile = "assets/" + filepath.Base(*inputFile) + "_compressed.nzip"
	}
	if *mode == "compress" {
		// Compress the file
		err := CompressFile(*inputFile, *outputFile)
		if err != nil {
			fmt.Println("Error compressing file:", err)
			os.Exit(1)
		}

		fmt.Printf("File successfully compressed to: %s\n", *outputFile)
		getStats(*inputFile, *outputFile)
	} else if *mode == "decompress" {
		err := DecompressFile(*inputFile, *outputFile)
		if err != nil {
			_ = fmt.Errorf("error while decompressing file %s", err)
		}
		fmt.Printf("File successfully decompressed to: %s\n", *outputFile)
	} else {
		fmt.Printf("Invalid mode %s", *mode)
	}
}

func getStats(inputFile string, outputFile string) {
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		fmt.Printf("Error retrieving file info: %v\n", err)
	}
	// Get file size
	inputFileSize := fileInfo.Size()
	fmt.Printf("The size of %s is %d bytes\n", inputFile, inputFileSize)
	fileInfo, err = os.Stat(outputFile)
	if err != nil {
		fmt.Printf("Error retrieving file info: %v\n", err)
	}
	// Get file size
	outputFileSize := fileInfo.Size()
	fmt.Printf("The size of %s is %d bytes\n", outputFile, outputFileSize)
	if outputFileSize > inputFileSize {
		fmt.Printf("Compressed file has %d more bytes", (outputFileSize - inputFileSize))
	} else {
		fmt.Printf("Compressed file has save %d bytes, in percent %f", (inputFileSize - outputFileSize), outputFileSize/inputFileSize*100)

	}
}
