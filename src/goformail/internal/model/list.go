package model

type List struct {
	Name       string   `json:"name"`
	Recipients []string `json:"recipients"`
}

type ListWithId struct {
	Id   int   `json:"id"`
	List *List `json:"list"`
}
