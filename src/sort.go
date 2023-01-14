package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

type ByRecord []Record

type Record struct {
	Key   []byte
	Value []byte
}

// func (a ByRecord) Len() int           { return len(a) }
// func (a ByRecord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a ByRecord) Less(i, j int) bool { return a[i].Key < a[j].Key }

func readInChunks(filename string, chunkSize int, keySize int, valueSize int) ByRecord {
	f, err := os.Open(filename) // Open is read only

	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	records := make([]Record, 0)
	buf := make([]byte, chunkSize)
	reader := bufio.NewReader(f)
	result := bytes.NewBuffer(nil)

	for {
		n, err := reader.Read(buf) // read up to len(buf) bytes, maybe lesser
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			break
		}
		currRecord := Record{Key: buf[:keySize], Value: buf[keySize:chunkSize]}
		records = append(records, currRecord)
		// fmt.Println(string(buf[:n])) // Ouptut bytes read
		result.Write(buf[:n]) // write to byte buffer
	}
	// fmt.Println("Complete message:", string(result.Bytes())) // full output

	return records
}

func writeBuffered(filename string, output []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close() // might cause issues -> https://www.joeshaw.org/dont-defer-close-on-writable-files/

	writer := bufio.NewWriter(f)
	// _, err = writer.Write(output)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for _, data := range output {
		_, err = writer.WriteString(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush() // flush the buffer to the file
}

func convertRecordsToBytes(records ByRecord) []string {
	output := make([]string, 0)
	for _, record := range records {
		output = append(output, string(record.Key[:]), string(record.Value[:]))
	}
	return output
}

func recordComparator(key1 []byte, key2 []byte) bool {
	compareValue := bytes.Compare(key1, key2)
	if compareValue == 1 {
		return true
	}
	return false
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	records := readInChunks(os.Args[0], 100, 10, 90)

	// sort.Slice(records, func(i, j int) bool {
	// 	return recordComparator(records[i].Key, records[j].Key)
	// })

	output := convertRecordsToBytes(records)

	writeBuffered(os.Args[2], output)

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
}
