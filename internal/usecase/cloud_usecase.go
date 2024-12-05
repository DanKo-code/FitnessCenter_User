package usecase

import "context"

type CloudUseCase interface {
	PutObject(ctx context.Context, object []byte, name string) (string, error)
	GetObjectByName(ctx context.Context, name string) ([]byte, error)
	DeleteObject(ctx context.Context, name string) error
}
