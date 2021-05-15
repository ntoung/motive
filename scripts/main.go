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
)

func main() {
	csvfile, err := os.Open("../docs/Sample Data/motive_api_sample_data_customer_reviews.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	//for {
	for i := 0; i < 10; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		MakeRequest(record[1])
		// fmt.Printf("Index: %s Text: %s \n", record[0], record[1])
	}
}

type ScorePostMessage struct {
	Correlation_id string `json:"correlation_id"`
	Domain         string `json:"domain"`
	Data_channel   string `json:"data_channel"`
	Models         []string
	Industry       string `json:"industry"`
	Prompt         string `json:"prompt"`
	Themes         map[string][]string
	Documents      []Document `json:"documents"`
}

type Document struct {
	Document_id string `json:"document_id"`
	Text        string `json:"text"`
}

// type Client struct {
// 	Transport     RoundTripper
// 	CheckRedirect func(req *Request, via []*Request) error
// 	Jar           CookieJar
// 	Timeout       time.Duration
// }

func MakeRequest(document string) {

	message := ScorePostMessage{
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

	// message := []byte(`{
	// 	"correlation_id": "string",
	// 	"domain": "other",
	// 	"data_channel": "other",
	// 	"models": [
	// 	  "sentiment",
	// 	  "emotion"
	// 	],
	// 	"industry": "education",
	// 	"prompt": "string",
	// 	"themes": {
	// 	  "property1": [
	// 		"string"
	// 	  ],
	// 	  "property2": [
	// 		"string"
	// 	  ]
	// 	},
	// 	"documents": [
	// 	  {
	// 		"document_id": "abc_123",
	// 		"text": "I am excited to be scoring this document."
	// 	  }
	// 	]
	//   }`)

	PostScores(message)

}

func PostScores(m ScorePostMessage) {
	apiUrl := "https://api-data.motivesoftware.com"
	resource := "/scores/"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api-data.motivesoftware.com/scores/"

	// Body
	body, err := json.Marshal(m)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body)) // URL-encoded payload
	r.Header.Add("X-API-Key", "AIzaSyCuFkUzhasM1OdOPVq-L0b8HQta7QTCJmU")
	r.Header.Add("Content-Type", "application/json")

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	//fmt.Println(formatrequest.FormatRequest(r))

	// Do and print response
	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
}
