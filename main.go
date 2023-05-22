package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type FileData struct {
	MD5    string `json:"MD5"`
	SHA1   string `json:"SHA1"`
	SHA256 string `json:"SHA256"`
}

var (
	version   string = "2.0.0"
	buildTime string = ""
	arch      string = ""
)

func main() {
	var isSplit bool
	var filename string
	var chunkSize int
	var verbose bool
	flag.BoolVar(&isSplit, "t", false, "Split the file into chunks")
	flag.StringVar(&filename, "f", "", "Filename")
	flag.IntVar(&chunkSize, "s", 10, "Chunk size in MB")
	flag.BoolVar(&verbose, "h", false, "Print file hash values")
	// parsing flags
	versionFlag := flag.Bool("v", false, "Show the app version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", version)
		fmt.Println("Build Time:", buildTime)
		fmt.Println("System Architecture:", arch)
		return
	}

	if filename == "" {
		fmt.Println("Please specify a filename")
		return
	}

	var err error
	if isSplit {
		err = splitFile(filename, chunkSize)
		if err != nil {
			fmt.Printf("Error splitting file: %v\n", err)
			return
		}
	} else if verbose {
		err = printFileHash(filename)
		if err != nil {
			fmt.Printf("Error sum file hash: %v\n", err)
			return
		}
	} else {
		err = recoverFile(filename, verbose)
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

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Calculate total number of chunks
	chunkSizeBytes := int64(chunkSize * 1024 * 1024)
	totalChunks := int(fileInfo.Size()/chunkSizeBytes) + 1

	// Create JSON file to hold file data
	jsonFile, err := os.Create(filename + ".json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Calculate file hashes
	md5Hash, sha1Hash, sha256Hash, err := getFileHashes(file)
	if err != nil {
		return err
	}

	// Write file data to JSON file
	fileData := FileData{MD5: md5Hash, SHA1: sha1Hash, SHA256: sha256Hash}
	fileDataJSON, err := json.Marshal(fileData)
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(fileDataJSON)
	if err != nil {
		return err
	}

	// Print file hash values if verbose mode is on
	fmt.Printf("MD5: %s\nSHA1: %s\nSHA256: %s\n", md5Hash, sha1Hash, sha256Hash)

	// Reset file pointer to start of file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// Split file into chunks
	for i := 0; i < totalChunks; i++ {
		// Create chunk file
		chunkFileName := fmt.Sprintf("%s.%d", filename, i)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		// Write chunk data to file
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

func recoverFile(filename string, verbose bool) error {
	// Open JSON file to get file data
	jsonFile, err := os.Open(filename + ".json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var fileData FileData
	err = json.NewDecoder(jsonFile).Decode(&fileData)
	if err != nil {
		return err
	}

	// Create new file for recovered data
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create hash functions
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	// Process each chunk of the file
	i := 0
	for {
		// Open next chunk file
		chunkFileName := fmt.Sprintf("%s.%d", filename, i)
		_, err := os.Stat(chunkFileName)
		if os.IsNotExist(err) {
			// No more chunks, verify hash values
			md5Value := fmt.Sprintf("%x", md5Hash.Sum(nil))
			sha1Value := fmt.Sprintf("%x", sha1Hash.Sum(nil))
			sha256Value := fmt.Sprintf("%x", sha256Hash.Sum(nil))

			if md5Value != fileData.MD5 || sha1Value != fileData.SHA1 || sha256Value != fileData.SHA256 {
				return fmt.Errorf("hash values do not match")
			}

			break
		}

		chunkFile, err := os.Open(chunkFileName)
		if err != nil {
			return err
		}
		defer chunkFile.Close()

		// Update hash values and write chunk data to file
		io.Copy(io.MultiWriter(file, md5Hash, sha1Hash, sha256Hash), chunkFile)

		i++
	}

	// Print verification result to console
	fmt.Println("File verification successful!")

	// Print file hash values if verbose mode is on
	if verbose {
		fmt.Printf("MD5: %s\nSHA1: %s\nSHA256: %s\n", fileData.MD5, fileData.SHA1, fileData.SHA256)
	}

	return nil
}

func printFileHash(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var md5Hash, sha1Hash, sha256Hash string
	md5Hash, sha1Hash, sha256Hash, err = getFileHashes(file)
	if err != nil {
		return err
	}

	fmt.Println("MD5:", md5Hash)
	fmt.Println("SHA-1", sha1Hash)
	fmt.Println("SHA-256", sha256Hash)

	return nil
}

func getFileHashes(file *os.File) (string, string, string, error) {
	// Create hash functions
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	// Process file data
	_, err := io.Copy(io.MultiWriter(md5Hash, sha1Hash, sha256Hash), file)
	if err != nil {
		return "", "", "", err
	}

	// Return hash values as strings
	md5Value := fmt.Sprintf("%x", md5Hash.Sum(nil))
	sha1Value := fmt.Sprintf("%x", sha1Hash.Sum(nil))
	sha256Value := fmt.Sprintf("%x", sha256Hash.Sum(nil))

	return md5Value, sha1Value, sha256Value, nil
}
