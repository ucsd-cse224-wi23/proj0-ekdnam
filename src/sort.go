package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"sort"
)

type ByRecord []Record

type Record struct {
	Key   []byte
	Value []byte
}

func (a ByRecord) Len() int           { return len(a) }
func (a ByRecord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRecord) Less(i, j int) bool { return string(a[i].Key) < string(a[j].Key) }

func readWholeFile(filename string, chunkSize int, keySize int, valueSize int) []Record {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	size := len(b) / chunkSize
	records := make([]Record, size)
	idx := 0
	// fmt.Printf("File has in total %d records\n", len(b)/100)
	storeIndex := 0
	for idx < len(b) {
		key := b[idx : idx+keySize]
		value := b[idx+keySize : idx+chunkSize]
		record := Record{Key: key, Value: value}
		records[storeIndex] = record
		storeIndex += 1
		idx += chunkSize
	}
	return records
}

func writeBuffered(filename string, output []string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close() // might cause issues -> https://www.joeshaw.org/dont-defer-close-on-writable-files/

	writer := bufio.NewWriter(f)

	for _, data := range output {
		_, err = writer.WriteString(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush() // flush the buffer to the file
}

func convertRecordsToString(records ByRecord) []string {
	output := make([]string, 0)
	for _, record := range records {
		output = append(output, string(record.Key[:]), string(record.Value[:]))
	}
	return output
}

func recordComparator(key1 []byte, key2 []byte) bool {
	compareValue := bytes.Compare(key1, key2)
	if compareValue == -1 {
		return true
	}
	return false
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}
	records := readWholeFile(os.Args[1], 100, 10, 90)

	sort.Sort(ByRecord(records))

	output := convertRecordsToString(records)

	writeBuffered(os.Args[2], output)

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])
}
