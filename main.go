package main

import (
	"encoding/json"
	"log"
	"net/http"

	// "database/sql"
	myDB "github.com/jyl/golang-TodoList/db"
	myType "github.com/jyl/golang-TodoList/type"
	_ "github.com/lib/pq"

	// myTypes "github.com/jyl/golang-TodoList/type"
	// "time"
	"html/template"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

const (
	templateDirectory = "./public/template/"
)

var wg sync.WaitGroup

var allFiles []string
var (
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
)

func populateTemplates() {
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

func deleteTask(w http.ResponseWriter, r *http.Request) {
	
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
	error := myDB.DeleteTask(taskID)
	if error != nil {
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

func editTaskFunc(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/edit/"))
	if err != nil {
		log.Fatal(err)
	}
	editTemplate.Execute(w, taskID)
}

func editTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	taskID,err := strconv.Atoi(r.FormValue("task-id"))
	if err != nil {
		log.Fatal(err)
	}
	updatedTitle := r.FormValue("task-title")
	updatedContent := r.FormValue("task-content")
	if updatedTitle != "" && updatedContent != "" {
		log.Println("passing info to DB ")
		error := myDB.EditTask(taskID, updatedTitle, updatedContent)
		if error != nil {
			log.Fatal("Error updating task")
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated task"))
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

func showAllTasksFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		context := myDB.GetTasks()

		// For each context.tasks, assign a goroutine to extract partial content
		for i := 0; i < len(context.Tasks); i++ {
			var task *myType.Task
			task = &context.Tasks[i]

			wg.Add(1)
			go extractTaskFunc(task)
		}
		wg.Wait()

		_, err := json.Marshal(context.Tasks)
		if err != nil {
			log.Fatal(err)
		}
		// w.Header().Set("Content-type", "application/json") // FIXME: this line causes Browser to interpret the template as plain text instead of rendering the plain text as html code
		w.WriteHeader(http.StatusOK)
		// w.Write(contextJson)
		homeTemplate.Execute(w, context)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func addTaskFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("enter addTaskFunc, method post?: ", r.Method)
	if r.Method == "POST" {
		title := r.FormValue("task-title")
		content := r.FormValue("task-content")
		log.Println(title)
		log.Println(content)
		if title != "" && content != "" {

			truth := myDB.AddTask(title, content)
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

func router(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		
		editTaskFuncAPI(w, r)
	} else if r.Method == "DELETE" {
		deleteTask(w,r)
	}

}

func main() {

	// http.HandleFunc("/complete/", CompleteTaskFunc)
	// http.HandleFunc("/delete/", DeleteTaskFunc)
	// http.HandleFunc("/deleted/", ShowTrashTaskFunc)
	// http.HandleFunc("/trash/", TrashTaskFunc)
	// http.HandleFunc("/completed/", ShowCompleteTasksFunc)
	// http.HandleFunc("/restore/", RestoreTaskFunc)
	// http.HandleFunc("/update/", UpdateTaskFunc)
	// http.HandleFunc("/search/", SearchTaskFunc)
	// http.HandleFunc("/login", GetLogin)
	// http.HandleFunc("/register", PostRegister)
	// http.HandleFunc("/admin", HandleAdmin)
	// http.HandleFunc("/add_user", PostAddUser)
	// http.HandleFunc("/change", PostChange)
	// http.HandleFunc("/logout", HandleLogout)
	populateTemplates()
	defer myDB.Close()
	http.HandleFunc("/task/", router)
	http.HandleFunc("/task", router)

	http.HandleFunc("/edit/", editTaskFunc)
	http.HandleFunc("/add/", addTaskFunc) // LEARNED: we registered a subtree named "/add/" and the form request is coming in with requesting handler for URL "/add" then ServeMux(which is an HTTP requst multiplexer matching URL of each incoming request against a list of registered patterns) redirects that request to the subtree root (adding the trailing slash). Note: HTTP rediect is called via GET method.
	http.HandleFunc("/add", addTaskFunc)
	// http.HandleFunc("/edit/", editTaskFunc)

	http.HandleFunc("/", showAllTasksFunc) // LEARNED: in Golang, url pattern "/" matches all paths not matched by other registered pattern - Example, if url "/gg" never exists in our app, Golang's ServeMux will match "/gg" to "/"

	// http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Print("running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // ListenAndServe always return non-nil error
}
