package module

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

type DocsService struct {
	service *docs.Service
}

func NewDocsService(ctx context.Context, b []byte) *DocsService {
	// サービスアカウントのクライアントを作成する
	srv, err := docs.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	return &DocsService{
		service: srv,
	}
}

func (d *DocsService) ReplaceAllText(documentId string, replaceMap map[string]string) {
	// 置換するテキストを設定するリクエスト
	requests := []*docs.Request{}
	for find, replace := range replaceMap {
		fmt.Printf("ReplaceAllText find: %#v\n", find)
		fmt.Printf("ReplaceAllText replace: %#v\n", replace)
		req := &docs.Request{
			ReplaceAllText: &docs.ReplaceAllTextRequest{
				ContainsText: &docs.SubstringMatchCriteria{
					Text: find,
				},
				ReplaceText: replace,
			},
		}
		requests = append(requests, req)
	}
	// リクエストをバッチで実行
	batchUpdateReq := &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
	_, err := d.service.Documents.BatchUpdate(documentId, batchUpdateReq).Do()
	if err != nil {
		log.Fatalf("ドキュメントのテキストを置換できませんでした: %v", err)
	}
	fmt.Println("テキストの置換が完了しました。")
}
