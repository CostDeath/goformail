package model

type List struct {
	Name       string   `json:"name"`
	Recipients []string `json:"recipients"`
}
