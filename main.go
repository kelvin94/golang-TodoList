package main

import (
	"encoding/json"
	"log"
	"net/http"

	// "database/sql"
	myDB "github.com/jyl/golang-TodoList/db"
	_ "github.com/lib/pq"

	// myTypes "github.com/jyl/golang-TodoList/type"
	// "time"
	"io/ioutil"
	"strings"
	"html/template"
)

const (
	templateDirectory = "./public/template/"
)

var err error
var allFiles []string

func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
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
	var homeTemplate *template.Template
	homeTemplate = templates.Lookup("home.html")
	// template.lookup() if not exists then Parse template file...
	if r.Method == "GET" {
		context := myDB.GetTasks()
		log.Println(context.Tasks)

		// FIXME: Decide whether build a rest api or template
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

func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	title := "random title"
	content := "random content"
	truth := myDB.AddTask(title, content)
	if truth != nil {
		log.Fatal("Error adding task")
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Added task"))
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
	defer myDB.Close()
	http.HandleFunc("/add/", AddTaskFunc)

	http.HandleFunc("/", ShowAllTasksFunc)

	// http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Print("running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // ListenAndServe always return non-nil error
}
