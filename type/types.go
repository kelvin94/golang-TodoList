package types



type Task struct {
    Id      int
    Title   string
    Content string
    Created string
    News       []News
}

//Context is the struct passed to templates
type Context struct {
    Tasks      []Task
    News       []News
}

type News struct {
	Title string `json:"title"`
    Url string `json:"url"`
    TaskId int `json:"taskId"`
}