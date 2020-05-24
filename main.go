package main

import(
	"net/http"
    "log"
    "encoding/json"
    // "database/sql"
    _ "github.com/lib/pq"
    myDB "github.com/jyl/golang-TodoList/db"
    // myTypes "github.com/jyl/golang-TodoList/type"
    // "time"
)




var err error


  

func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
        context := myDB.GetTasks()
        contextJson, err := json.Marshal(context.Tasks)
        if err != nil {
            log.Fatal(err)
        }
        w.Header().Set("Content-type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(contextJson)
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
    w.Write([]byte("Added task"))
}





func main() {
    
    // http.HandleFunc("/complete/", CompleteTaskFunc)
    // http.HandleFunc("/delete/", DeleteTaskFunc)
    // http.HandleFunc("/deleted/", ShowTrashTaskFunc)
    // http.HandleFunc("/trash/", TrashTaskFunc)
    // http.HandleFunc("/edit/", EditTaskFunc)
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
    log.Fatal(http.ListenAndServe(":8080", nil))
}