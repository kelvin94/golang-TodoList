package db

import (
    "database/sql"
    "errors"
	myTypes "github.com/jyl/golang-TodoList/type"
	_ "github.com/lib/pq"
	"log"
    "time"
    "strings"
)

/**
*
docker run -p 5432:5432 --name postgres_db -v postgres-volume:/var/lib/postgresql/data -d postgres
*/
type TaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(DB *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: DB,
	}
}

/***
sql.DB object performs tasks for you behind the scenes:
    1. It opens and closes connections to the actual underlying database, via the driver.
    2.It manages a pool of connections as needed, which may be a variety of things as mentioned.
*/
// var database *sql.DB
var err error

func (repo TaskRepository) DeleteTask(taskId int) error {
	deleteTaskQuery := "delete from task where id=$1;"
    
	log.Println("DeleteTask...  taskId", taskId)
	
    _, err := repo.db.Exec(deleteTaskQuery, taskId)
    // create a slice for the errors
    var errstrings []string
    var returnedError error
	if err != nil {
        log.Fatal(err)
        errstrings = append(errstrings, err.Error())
        returnedError = errors.New(strings.Join(errstrings, "\n"))
    } else {
        log.Println("delete DB success")
    }
   

	return returnedError
}

func (repo TaskRepository) GetTasks() myTypes.Context {
	var task []myTypes.Task
	var context myTypes.Context
	var TaskID int
	var TaskTitle string
	var TaskContent string
	var TaskCreated time.Time
	var getTasksql string

	getTasksql = "select id, title, content, created_date from task order by created_date asc;"

	rows, err := repo.db.Query(getTasksql)
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
		log.Println("ts: ", TaskCreated)
		a := myTypes.Task{Id: TaskID, Title: TaskTitle, Content: TaskContent,
			Created: TaskCreated.String()}
		task = append(task, a)
	}
	context = myTypes.Context{Tasks: task}
	return context
}

func (repo TaskRepository) EditTask(taskId int, title string, content string) error {
	query := "update task set title = $1, content = $2 where id = $3;"
	restoreSQL, err := repo.db.Prepare(query)
	log.Println("Editing Task... sql is done prepared query: " + query)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := repo.db.Begin()
	_, err = tx.Stmt(restoreSQL).Exec(title, content, taskId)
	if err != nil {

		log.Fatal(err)
		tx.Rollback()
	} else {
		log.Println("Update DB success")
		tx.Commit()
	}
	return err

}

func (repo TaskRepository) AddTask(title string, content string) error {

	query := "insert into task(title, content, created_date, last_modified_at) values($1,$2,now(), now())"
	restoreSQL, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := repo.db.Begin()
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

func (repo TaskRepository) DeleteNews(taskId int, title string) error {
	query := "delete from news where title = $1 and taskid=$2;"
	restoreSQL, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := repo.db.Begin()
	_, err = tx.Stmt(restoreSQL).Exec(title, taskId)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
	} else {
		log.Println("Delete DB success")
		tx.Commit()
	}
	return err
}

func (repo TaskRepository) AddNews(taskId int, title string, url string) error {
	query := "insert into news(title, url, taskId) values($1, $2, $3);"
	restoreSQL, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := repo.db.Begin()
	_, err = tx.Stmt(restoreSQL).Exec(title, url, taskId)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
	} else {
		log.Println("Insert DB success")
		tx.Commit()
	}
	return err
}

func (repo TaskRepository) GetNewsByTaskId(taskId int) []myTypes.News {

	query := "select title, url from news where taskid=$1"

	log.Println("taskId", taskId)

	rows, err := repo.db.Query(query, taskId)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var Url string
	var Title string
	var news []myTypes.News
	log.Println("rows", rows)
	if rows != nil {

		for rows.Next() {
			err := rows.Scan(&Title, &Url) // rows.Scan() scans result sets one row at a time and read the columns in each row into variables
			if err != nil {
				log.Fatal(err)
			}

			a := myTypes.News{TaskId: taskId, Title: Title, Url: Url}
			news = append(news, a)
		}
	}
	return news
}

func (repo TaskRepository) Close() {
	repo.db.Close()
	log.Println("DB connection is closed")
}
