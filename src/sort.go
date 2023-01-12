package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type ByRecord []Record

type Record struct {
	Key   byte
	Value byte
}

func (a ByRecord) Len() int           { return len(a) }
func (a ByRecord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRecord) Less(i, j int) bool { return a[i].Key < a[j].Key }

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	const BufferSize = 100
	const KeySize = 10
	const ValueSize = 90

	file, err := os.Open(os.Args[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)
	bufferKey := make([]byte, KeySize)
	bufferValue := make([]byte, ValueSize)

	records := make([]Record, 0)

	for {
		bytesKey, errKey := file.Read(bufferKey)
		bytesValue, errValue := file.Read(bufferValue)

		if errKey != nil {
			if errKey != io.EOF {
				fmt.Println(errKey)
			}

			break
		}
		if errValue != nil {
			if errValue != io.EOF {
				fmt.Println(errValue)
			}

			break
		}

		records = append(records, Record{bytesKey, bytesValue})
	}
	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
}
