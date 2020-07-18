package main

type Post struct {
	ID    string `json:"id" orm:"pk"`
	Title string `json:"title" orm:"type:varchar(255)"`
}

type User struct {
	ID   string `json:"id" orm:"pk"`
	Name string `json:"name"`
}

func main() {

}
