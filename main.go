package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var isSplit bool
	var filename string
	var chunkSize int
	flag.BoolVar(&isSplit, "t", false, "Split the file into chunks")
	flag.StringVar(&filename, "f", "", "Filename")
	flag.IntVar(&chunkSize, "s", 10, "Chunk size in MB")
	flag.Parse()

	if filename == "" {
		fmt.Println("Please specify a filename")
		return
	}

	if isSplit {
		err := splitFile(filename, chunkSize)
		if err != nil {
			fmt.Printf("Error splitting file: %v\n", err)
			return
		}
	} else {
		err := recoverFile(filename)
		if err != nil {
			fmt.Printf("Error recovering file: %v\n", err)
			return
		}
	}
}

func splitFile(filename string, chunkSize int) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	chunkSizeBytes := int64(chunkSize * 1024 * 1024)
	totalChunks := int(fileInfo.Size()/chunkSizeBytes) + 1

	for i := 0; i < totalChunks; i++ {
		chunkFileName := fmt.Sprintf("%s.%d", filename, i)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		buf := make([]byte, 4*1024) // 4KB buffer
		var bytesWritten int64

		for bytesWritten < chunkSizeBytes {
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			n, err = chunkFile.Write(buf[:n])
			if err != nil {
				return err
			}

			bytesWritten += int64(n)
		}
	}

	return nil
}

func recoverFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	i := 0
	for {
		chunkFileName := fmt.Sprintf("%s.%d", filename, i)
		_, err := os.Stat(chunkFileName)
		if os.IsNotExist(err) {
			// No more chunks
			break
		}

		chunkFile, err := os.Open(chunkFileName)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		_, err = io.Copy(file, chunkFile)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}
