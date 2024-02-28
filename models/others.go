package models

type Fields struct {
	Name  string `json:"name" bson:"name"`
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}

type Message struct {
	Content   string        `json:"content" bson:"content"`
	Status    MessageStatus `bson:"status"`
	Timestamp int64         `json:"timestamp" bson:"timestamp"`
	Sender    string        `json:"sender" bson:"sender"`
}

type MessageStatus struct {
	Sent     bool `json:"sent" bson:"sent"`
	Received bool `json:"received" bson:"received"`
}

// / Whatsapp
type MessagePayload struct {
	Type    string
	Content string
}

type Row struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
