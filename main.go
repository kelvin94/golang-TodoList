package main

import (
	"log"
	"net/http"

	"database/sql"
	"github.com/jyl/golang-TodoList/db"
	"github.com/jyl/golang-TodoList/app"
	"github.com/jyl/golang-TodoList/api"
	// "github.com/jyl/golang-TodoList/structs"

	_ "github.com/lib/pq"

	// myTypes "github.com/jyl/golang-TodoList/type"
	// "time"
	
)


func editHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fdsfds"))
		h.ServeHTTP(w, r)
	  })
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
	connStr := "host=localhost user=golang password=golang dbname=golang sslmode=disable"
    database, err := sql.Open("postgres", connStr)
	var taskRepository *db.TaskRepository = db.NewPostgresTaskRepository(database)
	
	if err != nil {
		log.Fatal(err)
    } else {
        log.Println("DBConnection success")
	}

	
	http.Handle("/api/task", &api.Api{Repo : taskRepository})
	

	
	var myAppHandler *app.MyApp = &app.MyApp{Repo : taskRepository}
	http.HandleFunc("/", myAppHandler.ShowAllTasksFunc) // LEARNED: in Golang, url pattern "/" matches all paths not matched by other registered pattern - Example, if url "/gg" never exists in our app, Golang's ServeMux will match "/gg" to "/"
	http.HandleFunc("/add/", myAppHandler.AddTaskFunc) // LEARNED: we registered a subtree named "/add/" and the form request is coming in with requesting handler for URL "/add" then ServeMux(which is an HTTP requst multiplexer matching URL of each incoming request against a list of registered patterns) redirects that request to the subtree root (adding the trailing slash). Note: HTTP rediect is called via GET method.
	http.HandleFunc("/add", myAppHandler.AddTaskFunc)
	http.HandleFunc("/edit/", myAppHandler.EditRouter)	
	http.HandleFunc("/delete", myAppHandler.DeleteTaskFunc)

	// http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Print("running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // ListenAndServe always return non-nil error
}
