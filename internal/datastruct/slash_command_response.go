package datastruct

type SlashCommandResponse struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

type Attachment struct {
	Text     string `json:"text"`
	ImageUrl string `json:"image_url"`
}
