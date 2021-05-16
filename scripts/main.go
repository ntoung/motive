package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type ScorePostMessage struct {
	Correlation_id string   `json:"correlation_id"`
	Domain         string   `json:"domain"`
	Data_channel   string   `json:"data_channel"`
	Models         []string `json:"models"`
	Industry       string   `json:"industry"`
	Prompt         string   `json:"prompt"`
	Themes         map[string][]string
	Documents      []Document `json:"documents"`
}

type Document struct {
	Document_id string `json:"document_id"`
	Text        string `json:"text"`
}

func MakePostScoreRequestPayload(document string) ScorePostMessage {
	return ScorePostMessage{
		Correlation_id: "string",
		Domain:         "other",
		Data_channel:   "other",
		Models: []string{
			"sentiment",
			"emotion",
		},
		Industry: "education",
		Prompt:   "string",
		Themes: map[string][]string{
			"property1": {"string"},
			"property2": {"string"},
		},
		Documents: []Document{
			{
				Document_id: "abc_123",
				Text:        document,
			},
		},
	}
}

func PostScores(document string, ch chan<- string) {
	payload := MakePostScoreRequestPayload(document)
	apiUrl := "https://api-data.motivesoftware.com"
	resource := "/scores/"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api-data.motivesoftware.com/scores/"

	// Body
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body))
	r.Header.Add("X-API-Key", "AIzaSyCuFkUzhasM1OdOPVq-L0b8HQta7QTCJmU")
	r.Header.Add("Content-Type", "application/json")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	//ch <- fmt.Sprint(string(requestDump))
	fmt.Println(string(requestDump))

	// Do and print response
	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
}

func main() {
	csvfile, err := os.Open("../docs/Sample Data/motive_api_sample_data_customer_reviews.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	start := time.Now()
	ch := make(chan string)
	r := csv.NewReader(csvfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		go PostScores(record[1], ch)
	}
	//fmt.Println(<-ch)
	fmt.Println("END")
	fmt.Println(time.Since(start))
}
