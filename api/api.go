package api

import(
	"log"
	"net/http"
	"strings"
	"encoding/json"
	myDB "github.com/jyl/golang-TodoList/db"
	myType "github.com/jyl/golang-TodoList/type"
    "fmt"
	"io/ioutil"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}


type Api struct {
	Repo *myDB.TaskRepository
}

func (api Api) ServeHTTP(w http.ResponseWriter, r *http.Request){
	// match the incoming URL with regular expression
	// GET URL: /api/task?title=someTitle
	// POST URL: /api/task
	// PUT URL: /api/task
	// Delete URL: /api/task
	hnDataChan := make(chan []byte)
	done := make(chan struct{},50)
	newsChan := make(chan myType.News,50)
	switch r.Method {
		case http.MethodGet:
			api.get(w, r, hnDataChan, done, newsChan)
		case http.MethodPost:
			api.post(w,r)
		case http.MethodDelete:
			api.delete(w,r)
		default:
			log.Fatal("A request with unexpected HTTP method: ", r)
		
	}
	
}

func (api Api) delete(w http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		// pointer so we can test for field absence
		Title *string `json:"title"`
		TaskId *string `json:"taskId"`
	}{}

	err := d.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Title == nil {
		http.Error(w, "missing field 'title' from JSON object", http.StatusBadRequest)
		return
	} else if t.TaskId == nil {
		http.Error(w, "missing field 'TaskId' from JSON object", http.StatusBadRequest)
		return
	}
	taskId, er := strconv.Atoi(*t.TaskId)
	if er != nil {
		log.Fatal(er)
	}
	e := api.Repo.DeleteNews(taskId, *t.Title)
	if e != nil {
		log.Fatal(er)
		http.Error(w, "DB insert news instance failed", http.StatusBadRequest)
		return
	}
	jData, err := json.Marshal("success")
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jData)
}


func (api Api) post(w http.ResponseWriter, req *http.Request) {

	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	t := struct {
		// pointer so we can test for field absence
		Title *string `json:"title"`
		Url *string `json:"url"`
		TaskId *string `json:"taskId"`
	}{}

	err := d.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t.Url == nil  {
		http.Error(w, "missing field 'Url' from JSON object", http.StatusBadRequest)
		return
	} else if t.Title == nil {
		http.Error(w, "missing field 'title' from JSON object", http.StatusBadRequest)
		return
	} else if t.TaskId == nil {
		http.Error(w, "missing field 'TaskId' from JSON object", http.StatusBadRequest)
		return
	}
	taskId, err := strconv.Atoi(*t.TaskId)
	if err != nil {
		log.Fatal(err)
	}
	er := api.Repo.AddNews(taskId, *t.Title, *t.Url)
	if er != nil {
		log.Fatal(er)
		http.Error(w, "DB insert news instance failed", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	
}

func (api Api) get(w http.ResponseWriter, r *http.Request, hnDataChan chan []byte,  done chan struct{}, newsChan chan myType.News) {
	title := r.URL.Query()["title"]
	hnApi := "https://hn.algolia.com/api/v1/search?query=" // Hacker News Api endpoint

	var str strings.Builder
	str.WriteString(hnApi)
	str.WriteString(title[0])
	

	runChannelListener(hnDataChan, done, newsChan, w)
	wg.Add(1)  
	go func(news chan myType.News, w http.ResponseWriter, done chan struct{}) { // this is the receiving channel, 下面有个例子能够优化接收方，让接收方find out that the channel is closed 
		defer wg.Done()
		var newsArr []myType.News
		for n := range news { // LEARNED: because we rely on for range loop to listen on the NewsChannel, if NewsChannel is never closed, then this loop will be running forever, so that we need to close NewsChannel in the processHNApiData function to notify the goroutine "go func(news chan News, w http.ResponseWriter, done chan struct{})" that "News channel got no more data to send, end the for...range loop and carry on your task"
			newsArr = append(newsArr, n)
		}

		jData, err := json.Marshal(newsArr)
		if err != nil {
			// handle error
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		 
		log.Println("sending signal to done channel")
		done <- struct{}{}

	}(newsChan, w, done)  

	// start 10 go routinues that run a select statement to receive messages from hnDataChannel and 
	// but for now, we only search hacker news for 1 title
	// for each title... run go getHNApiData()...
	getHNApiData(hnDataChan, str.String())
	wg.Wait()
}

func runChannelListener(hnDataChan chan []byte,  done chan struct{}, newsChan chan myType.News, w http.ResponseWriter) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
				case msg := <- hnDataChan:
					processHNApiData(msg, done, newsChan)
				// case msg := <- newsChan:
				// 	log.Println("news channel received")
				// 	writeBackToClient(msg, done, w)
				
				case <- done:
					close(hnDataChan)
					
					close(done)
					return
			}

		}
	}()
}

func getHNApiData( hnDataChan chan []byte, url string) {
	
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		hnDataChan <- data;	
	}
	
	
}


type Data struct {
    MyKey []interface{} `json:"hits"`
}

func processHNApiData( msg []byte,  done chan struct{}, newsChan chan myType.News) {
	
	
	var data Data

	err := json.Unmarshal(msg, &data) 
	
	if err != nil {
		log.Fatal(err)
	}
	
	for i, v := range data.MyKey {
		switch x := v.(type) {
			case map[string]interface{}:
				if x["title"] != "" && x["title"] != nil && x["url"] != "" && x["url"] != nil {
					t := fmt.Sprintf("%v", x["title"])
					u := fmt.Sprintf("%v", x["url"])				
	
					if t != "" && len(t) > 0 && u != "" && len(u) > 0 {	
					
						n := myType.News{Title : t, Url: u}
						newsChan <- n
					}
				}
				
			default:
				log.Printf("%d. Unexpected value of type %T\n", i, v)
		}
	}
	close(newsChan) // closing newsChannel when finishing sending data to the NewsChannel
}

