package app

import (
	"encoding/json"
	"log"
	"net/http"

	myDB "github.com/jyl/golang-TodoList/db"
	myType "github.com/jyl/golang-TodoList/type"

	// myTypes "github.com/jyl/golang-TodoList/type"
	"strconv"
	"html/template"
	"io/ioutil"
	"strings"
	"sync"
)

const (
	templateDirectory = "./public/template/"
)
 
var (
	allFiles 			[]string
	homeTemplate        *template.Template
	addTaskFormTemplate *template.Template
	deletedTemplate     *template.Template
	completedTemplate   *template.Template
	loginTemplate       *template.Template
	editTemplate        *template.Template
	searchTemplate      *template.Template
	templates           *template.Template
	message             string
	//message will store the message to be shown as notification
	err error
	wg sync.WaitGroup

)

type MyApp struct {
    Repo *myDB.TaskRepository

}

func init() {
	files, err := ioutil.ReadDir(templateDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templateDirectory+filename)
		}
	}
	templates := template.Must(template.ParseFiles(allFiles...)) // Must is a helper that wraps a call to a function returning (*Template, error) and panics if the error is non-nil; Also, Must() function used to verify that a template is valid during parsing.
	homeTemplate = templates.Lookup("home.html")
	editTemplate = templates.Lookup("EditTaskForm.html")
	addTaskFormTemplate = templates.Lookup("AddTaskForm.html")
}

func (myApp MyApp) DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {

	// LEARNED: You've sent your request with a JSON body, but ParseForm on an *http.Request does not handle JSON. You need to read the body of the request and parse it as JSON, or don't send your body as JSON.
	var body = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	taskID, err := strconv.Atoi(body["task_id"])
	if err != nil {
		log.Fatal(err)
	}
	error := myApp.Repo.DeleteTask(taskID)
	if error != nil {
		log.Println("error not null")
		log.Fatal(error)
	}

	// LEARNED: 因为前端跟后端握手时候决定了是cappliation/json content-type所以前端的ajax success event也在等返回一个json as data to trigger the "success" callback in the ajax call.
	jData, err := json.Marshal("true")
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func (myApp MyApp) EditRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			myApp.EditTaskFunc(w, r)
		case http.MethodGet:
			myApp.EditPage(w,r)
		default:
			log.Fatal("Edit router receives unknown request with unexpected method")
	}
}

func (myApp MyApp) EditPage(w http.ResponseWriter, r *http.Request) {

	taskID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/edit/"))
	if err != nil {
		log.Fatal(err)
	}
	editTemplate.Execute(w, taskID)
}

func (myApp MyApp) EditTaskFunc(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.FormValue("task-id"))
	if err != nil {
		log.Fatal(err)
	}
	updatedTitle := r.FormValue("task-title")
	updatedContent := r.FormValue("task-content")
	if updatedTitle != "" && updatedContent != "" {
		error := myApp.Repo.EditTask(taskID, updatedTitle, updatedContent)
		if error != nil {
			log.Fatal("Error updating task")
		}
		http.Redirect(w, r, "/", http.StatusFound)

		// w.Header().Set("Content-type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("Updated task"))
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func extractTaskFunc(task *myType.Task) {
	if len(task.Content) > 10 {
		task.Content = task.Content[:10]
	}
	wg.Done()
}

func (myApp MyApp) ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		context := myApp.Repo.GetTasks()
		// For each context.tasks, assign a goroutine to extract partial content
		for i := 0; i < len(context.Tasks); i++ {
			var task *myType.Task
			task = &context.Tasks[i]
			news := myApp.Repo.GetNewsByTaskId(task.Id)

			task.News = news
			
			wg.Add(1)
			go extractTaskFunc(task)
		}
		wg.Wait()

		_, err := json.Marshal(context.Tasks)
		if err != nil {
			log.Fatal(err)
		}
		// w.Header().Set("Content-type", "application/json") // FIXME: this line causes Browser to interpret the template as plain text instead of rendering the plain text as html code
		log.Println("context", context)
		w.WriteHeader(http.StatusOK)
		// w.Write(contextJson)
		homeTemplate.Execute(w, context)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func (myApp MyApp) AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("enter addTaskFunc, method post?: ", r.Method)
	if r.Method == "POST" {
		title := r.FormValue("task-title")
		content := r.FormValue("task-content")
		log.Println(title)
		log.Println(content)
		if title != "" && content != "" {

			truth := myApp.Repo.AddTask(title, content)
			if truth != nil {
				log.Fatal("Error adding task")
			}
			http.Redirect(w, r, "/", http.StatusFound)

		} else {

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to add task"))
		}
	} else if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		addTaskFormTemplate.Execute(w, nil)
	}
}