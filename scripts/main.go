package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func postScores(document string, ch chan<- string) postScoresResponse {
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
	req, _ := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(body))
	req.Header.Add("X-API-Key", "AIzaSyCuFkUzhasM1OdOPVq-L0b8HQta7QTCJmU")
	req.Header.Add("Content-Type", "application/json")

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// // Do and print response
	// resp, _ := client.Do(req)
	// fmt.Println(resp.Status)

	// var result map[string]interface{}
	// json.NewDecoder(resp.Body).Decode(&result)
	// log.Println(result)
	// return result

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	postScoresRes := postScoresResponse{}
	jsonErr := json.Unmarshal(body, &postScoresRes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return postScoresRes 
}

type ScoreStatus string

const (
	Processing ScoreStatus = "PROCESSING"
	Enqueued   ScoreStatus = "ENQUEUED"
	Done       ScoreStatus = "DONE"
)

type postScoresResponse struct {
	Job_id string `json:"job_id"`
}

type getScoresResponse struct {
	Job_Id string      `json:"job_id"`
	Status ScoreStatus
}

func getScores(jobId string) getScoresResponse {
	apiUrl := "https://api-data.motivesoftware.com"
	resource := "scores"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%s/%s/", u.String(), jobId) // "https://api-data.motivesoftware.com/scores"

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, urlStr, nil)
	req.Header.Add("X-API-Key", "AIzaSyCuFkUzhasM1OdOPVq-L0b8HQta7QTCJmU")

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	getScoresRes := getScoresResponse{}
	jsonErr := json.Unmarshal(body, &getScoresRes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	log.Println("=============")
	log.Println(getScoresRes)
	log.Println(getScoresRes)

	return getScoresRes 
}

// func concurrentPost() {
// 	csvfile, err := os.Open("../docs/Sample Data/motive_api_sample_data_customer_reviews.csv")
// 	if err != nil {
// 		log.Fatalln("Couldn't open the csv file", err)
// 	}

// 	start := time.Now()
// 	ch := make(chan string)
// 	r := csv.NewReader(csvfile)
// 	//for k := 1; k <= 10; k++ {
// 	for {
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		go postScores(record[1], ch)
// 		go postScores(record[1], ch)
// 		go postScores(record[1], ch)
// 		go postScores(record[1], ch)
// 		go postScores(record[1], ch)
// 	}
// 	//}
// 	//fmt.Println(<-ch)
// 	fmt.Println("END")
// 	fmt.Println(time.Since(start))
// }

// func worker(job_id string, jobs <- chan int, results chan<- int) {
// 	for j := range jobs {
// 		scoreResponse := getScores(job_id)
// 		switch scoreResponse.Status {
// 			case Processing: 
// 				log.Printf("PROCESSING\n")
// 			case Enqueued: 
// 				log.Printf("ENQUEUED\n")
// 			case Done: 
// 				log.Printf("DONE\n")
// 		}
// 		results <- j * 2
// 	}
// }

func postAndPollForResponse() {
	csvfile, err := os.Open("../docs/Sample Data/motive_api_sample_data_customer_reviews.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	start := time.Now()
	ch := make(chan string)
	r := csv.NewReader(csvfile)
	for k := 1; k <= 1; k++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		result := postScores(record[1], ch)

		log.Println("polling status")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		// Having cancelFunc here gives control if you want to cancel the poll before deadline
		// but for this instance, I'm just gonna defer cancel it
		defer cancelFunc()
		status, err := PollStatus(ctx, result.Job_id, 1*time.Second)
		if err != nil {
			log.Fatalf("PollStatus error: %+v", err)
		}

		log.Print(status)

		//getScores(result.Job_id)

		
		// const numJobs = 5
		// jobs := make(chan int, numJobs)
		// results := make(chan int, numJobs)
		
		// for w := 1; w <= 3; w++ {
		// 	go worker(result.Job_id, jobs, results)
		// }

		// for j := 1; j <= numJobs; j++ {
		// 	jobs <- j
		// }

		// close(jobs)

		// for a := 1; a <= numJobs; a++ {
		// 	<-results
		// }


		// done := make(chan bool)
		// go func() {
		// 	for {
		// 		select {
		// 		case <- done:
		// 			scoreResponse := getScores(result.Job_id)
		// 			switch scoreResponse.Status {
		// 				case Processing: 
		// 					log.Printf("%d PROCESSING\n", i)
		// 				case Enqueued: 
		// 					log.Printf("%d ENQUEUED\n", i)
		// 				case Done: 
		// 					log.Printf("%d DONE\n", i)
		// 			}
		// 		}
		// 	}
		// }()	
		
	}
	//}
	//fmt.Println(<-ch)
	fmt.Println("END")
	fmt.Println(time.Since(start))
}

// PollStatus will keep 'pinging' the status API until timeout is reached or status returned is either successful or in error.
func PollStatus(ctx context.Context, id string, pollInterval time.Duration) (*getScoresResponse, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	tickerCounter := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case tick := <-ticker.C:
			res := getScores(id)
			tickerCounter++
			log.Printf("tick #%d at '%ss', %s", tickerCounter, tick.Format("05"), res.Status)
			if res.Status == Done {
				return &res, nil
			}

		}
	}

	//return nil, errors.New("unable to get status")
}

func main() {
	//concurrentPost()
	postAndPollForResponse()
}
