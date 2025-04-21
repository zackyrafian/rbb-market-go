package websocket

type Message struct {
	From    int    `json:"from"`
	To      int    `json:"to"`
	Content string `json:"content"`
}
