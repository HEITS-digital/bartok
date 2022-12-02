package repository

import (
	"cloud.google.com/go/firestore"
	"context"
)

type DAO interface {
	NewWatercoolerQuery() WatercoolerQuery
}
type dao struct {
	client *firestore.Client
}

func NewDAO(f *firestore.Client) DAO {
	return &dao{f}
}
func NewFirestoreClient(projectId string) (*firestore.Client, error) {

	ctx := context.Background()
	FirestoreClient, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return FirestoreClient, nil
}

func (d *dao) NewWatercoolerQuery() WatercoolerQuery {
	return &watercoolerQuery{client: d.client}
}
