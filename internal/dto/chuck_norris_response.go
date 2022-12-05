package dto

type ChuckNorrisResponse struct {
	Value      string   `json:"value"`
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Url        string   `json:"url"`
}
