package models


type Message struct {
	Username string `json:"username"` 
	Action   string `json:"action"`   
	Content  string `json:"content"`  
}
