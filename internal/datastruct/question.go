package datastruct

type Question struct {
	Id       string
	English  string `firestore:"question_en"`
	Romanian string `firestore:"question_ro"`
	IsRead   bool   `firestore:"is_read"`
}
