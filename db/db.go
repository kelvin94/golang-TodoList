package db

import(
	"net/http"
    "log"
    "database/sql"
	_ "github.com/lib/pq"
	"github.com/jyl/Tasks/type"
    "time"
)
/**
*
docker run -p 5432:5432 --name postgres_db -v postgres-volume:/var/lib/postgresql/data -d postgres
*/
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    dbname   = "postgres"
)
/***
sql.DB object performs tasks for you behind the scenes:
    1. It opens and closes connections to the actual underlying database, via the driver.
    2.It manages a pool of connections as needed, which may be a variety of things as mentioned.
*/
var database *sql.DB 
var err error

func GetTasks() Context {
    var task []Task
    var context Context
    var TaskID int
    var TaskTitle string
    var TaskContent string
    var TaskCreated time.Time
    var getTasksql string

    getTasksql = "select id, title, content, created_date from task;"

    rows, err := database.Query(getTasksql)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&TaskID, &TaskTitle, &TaskContent, &TaskCreated) // rows.Scan() scans result sets one row at a time and read the columns in each row into variables
        if err != nil {
            log.Fatal(err)
        }
        TaskCreated = TaskCreated.Local()
        a := Task{Id: TaskID, Title: TaskTitle, Content: TaskContent,
                    Created: TaskCreated.Format(time.UnixDate)[0:20]}
        task = append(task, a)
    }
    context = Context{Tasks: task}
    return context
}

func AddTask(title string, content string) error {
    query:="insert into task(title, content, created_date, last_modified_at) values(?,?,datetime(), datetime())"
    restoreSQL, err := database.Prepare(query)
    if err != nil {
        log.Fatal(err)
    }
    tx, err := database.Begin()
    _, err = tx.Stmt(restoreSQL).Exec(title, content)
    if err != nil {
        log.Fatal(err)
        tx.Rollback()
    } else {
        log.Println("Insert DB success")
        tx.Commit()
    }
    return err
}

func init() {
    
    connStr := "host=142.93.159.74 user=postgres password=ggininder87 dbname=golang sslmode=disable"
    database, err = sql.Open("postgres", connStr)
    
	if err != nil {
		log.Fatal(err)
    } else {
        log.Println("DBConnection success")
    }
    
    
}