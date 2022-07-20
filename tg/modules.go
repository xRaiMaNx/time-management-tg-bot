package tg

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}
