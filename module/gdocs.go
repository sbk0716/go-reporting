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

// サービスアカウントのクライアントを作成する
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

// 対象のドキュメントのテキストを置換する
func (d *DocsService) ReplaceAllText(documentId string, replacements map[string]string) {
	// 置換するテキストを設定するリクエスト
	requests := []*docs.Request{}
	for searchStr, replaceStr := range replacements {
		fmt.Printf("ReplaceAllText searchStr: %#v\n", searchStr)
		fmt.Printf("ReplaceAllText replaceStr: %#v\n", replaceStr)
		req := &docs.Request{
			ReplaceAllText: &docs.ReplaceAllTextRequest{
				ContainsText: &docs.SubstringMatchCriteria{
					Text: searchStr,
				},
				ReplaceText: replaceStr,
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
	fmt.Println("ドキュメントのテキストの置換が完了しました。")
}
