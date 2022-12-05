package repository

import (
	"bartok/internal/datastruct"
	"cloud.google.com/go/firestore"
	"context"
)

type WatercoolerQuery interface {
	GetNextUnreadQuestion() (*datastruct.Question, error)
	UpdateQuestion(*datastruct.Question) error
}

const questionsCollection = "questions"

type watercoolerQuery struct {
	client *firestore.Client
}

func (w *watercoolerQuery) UpdateQuestion(q *datastruct.Question) error {
	ctx := context.Background()
	_, err := w.client.Collection(questionsCollection).Doc(q.Id).Set(ctx, q)
	return err
}
func (w *watercoolerQuery) GetNextUnreadQuestion() (*datastruct.Question, error) {
	ctx := context.Background()

	questions := w.client.Collection(questionsCollection)
	q, err := questions.Where("is_read", "==", false).Limit(1).Documents(ctx).Next()
	if err != nil {
		return nil, err
	}
	var nextQuestion datastruct.Question
	if err := q.DataTo(&nextQuestion); err != nil {
		return nil, err
	}
	nextQuestion.Id = q.Ref.ID
	return &nextQuestion, nil
}
