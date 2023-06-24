package module

import (
	"context"
	"io/ioutil"
	"log"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

func NewDocsService(ctx context.Context, opts ...option.ClientOption) (*docs.Service, error) {
	// サービスアカウントの秘密鍵を読み込む
	b, err := ioutil.ReadFile("secret.json")
	if err != nil {
		log.Fatalf("秘密鍵ファイルを読み込めませんでした: %v", err)
	}

	// サービスアカウントのクライアントを作成する
	srv, err := docs.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	return srv, err
}
