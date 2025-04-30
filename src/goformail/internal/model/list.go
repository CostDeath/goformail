package model

type ListRequest struct {
	Name            string   `json:"name"`
	Recipients      []string `json:"recipients"`
	Mods            []int64  `json:"mods"`
	ApprovedSenders []string `json:"approved_senders"`
}

type ListResponse struct {
	Id              int      `json:"id"`
	Name            string   `json:"name"`
	Recipients      []string `json:"recipients"`
	Mods            []int64  `json:"mods"`
	ApprovedSenders []string `json:"approved_senders"`
}

type ListOverrides struct {
	Recipients      bool
	Mods            bool
	ApprovedSenders bool
}
