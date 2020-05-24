package main

import(
	"net/http"
    "log"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/jyl/Tasks/db"
    "github.com/jyl/Tasks/type"
    "time"
)




var err error


  

func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {
        context := GetTasks()
        w.Write([]byte(context.Tasks[0].Title))
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	
}



func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
    title := "random title"
    content := "random content"
    truth := AddTask(title, content)
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
    defer database.Close()
    http.HandleFunc("/add/", AddTaskFunc)
    http.HandleFunc("/", ShowAllTasksFunc)

    // http.Handle("/static/", http.FileServer(http.Dir("public")))
    log.Print("running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}