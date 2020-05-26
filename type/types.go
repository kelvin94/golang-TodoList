package types



type Task struct {
    Id      int
    Title   string
    Content string
    Created string
}

//Context is the struct passed to templates
type Context struct {
    Tasks      []Task
    Navigation string
    Search     string
    Message    string
}