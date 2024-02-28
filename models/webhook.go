package models

type WebhookPayload struct {
	Entry []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MetaData struct {
					PhoneNumberId string `json:"phone_number_id"`
				}
				Contacts []struct {
					WAID    string `json:"wa_id"`
					Profile struct {
						Name string `json:"name"`
					}
				} `json:"contacts"`
				Messages []struct {
					Type        string `json:"type"`
					Interactive struct {
						Type        string `json:"type"`
						ButtonReply struct {
							Id    string `json:"id"`
							Title string `json:"title"`
						} `json:"button_reply"`
						ListReply struct {
							Id    string `json:"id"`
							Title string `json:"title"`
						} `json:"list_reply"`
					} `json:"interactive"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
					Audio struct {
						Id    string `json:"id"`
						Mime  string `json:"mime_type"`
						Sha   string `json:"sha256"`
						Voice bool   `json:"voice"`
					} `json:"audio"`
				} `json:"messages"`
				Statuses []struct {
					Status string `json:"status"`
				} `json:"statuses"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}
