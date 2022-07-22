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

type BotMsg struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard []KeyboardButton `json:"keyboard"`
}

type BotMsgKeyboardMarkup struct {
	ChatID      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}
