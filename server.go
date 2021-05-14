package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	csvfile, err := os.Open("./docs/Sample Data/motive_api_sample_data_customer_reviews.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)

	for i := 0; i < 10; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Index: %s Text: %s", record[0], record[1])
	}
}

// func NewReader(r io.Reader) *Reader {
// 	return &Reader{
// 		Comma: ',',
// 		r:     bufio.NewReader(r),
// 	}
// }
