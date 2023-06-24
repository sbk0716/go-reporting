package module

import (
	"context"
	"log"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type DocsService struct {
	service *docs.Service
}

func NewDocsService(ctx context.Context, b []byte) (*docs.Service, error) {
	// サービスアカウントのクライアントを作成する
	srv, err := docs.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	return srv, err
}
