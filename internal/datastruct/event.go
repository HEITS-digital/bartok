package datastruct

type EmployeeEventType string

const (
	Birthday    EmployeeEventType = "birthday"
	Anniversary                   = "anniversary"
)

type EmployeeEvent struct {
	Id       string            `json:"id"`
	Type     EmployeeEventType `json:"type"`
	Employee string            `json:"employee"`
	Date     string            `json:"date"`
	Text     string            `json:"text"`
	IsSent   bool              `json:"is_sent"`
}
