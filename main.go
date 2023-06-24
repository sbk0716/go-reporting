package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// サービスアカウントの秘密鍵を読み込む
	b, err := ioutil.ReadFile("secret.json")
	if err != nil {
		log.Fatalf("秘密鍵ファイルを読み込めませんでした: %v", err)
	}

	// サービスアカウントのクライアントを作成する
	docSrv, err := docs.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	driveSrv, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}

	// 環境変数からファイルIDを取得
	sourceDocId := os.Getenv("DOC_ID")
	if sourceDocId == "" {
		// コピー元のドキュメントID
		// https://docs.google.com/document/d/1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg/edit
		sourceDocId = "1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg"
	}

	fullName := os.Getenv("FULL_NAME")
	if fullName == "" {
		fullName = "山田 太郎"
	}
	email := os.Getenv("EMAIL")
	if email == "" {
		email = "taro.yamada@test.com"
	}

	// タイムスタンプを取得（現在時刻をJSTに変換）
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp := time.Now().In(jst).Format("2006-01-02-15-04-05")

	// 複製先のGoogleドキュメントのタイトル
	// newDocumentTitle := "Copy-of-Document"
	newDocumentTitle := fmt.Sprintf("%s_Copy-of-Document", timestamp)

	// 複製リクエストを作成
	copyRequest := &drive.File{
		Title: newDocumentTitle,
	}

	// Googleドキュメントを複製
	copiedDocument, err := driveSrv.Files.Copy(sourceDocId, copyRequest).Do()
	if err != nil {
		log.Fatalf("Googleドキュメントの複製に失敗しました: %v", err)
	}
	// 複製先のGoogleドキュメントのIDを出力
	fmt.Printf("Googleドキュメントの複製が完了しました。複製先のドキュメントID: %s\n", copiedDocument.Id)
	r, err := driveSrv.Files.List().
		Fields("*").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	// ファイル一覧を表示する
	fmt.Println("ファイル一覧:")
	if len(r.Items) > 0 {
		for _, file := range r.Items {
			fmt.Printf("ファイル名: %s (ID: %s)\n", file.Title, file.Id)
		}
	} else {
		fmt.Println("ファイルが見つかりませんでした。")
	}
	copyDocId := copiedDocument.Id

	// 置換対象の文字列と置換後の文字列のマップ
	replaceMap := map[string]string{
		"${fullName}": fullName,
		"${email}":    email,
	}

	// 置換するテキストを設定するリクエスト
	requests := []*docs.Request{}
	for find, replace := range replaceMap {
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
	_, err = docSrv.Documents.BatchUpdate(copyDocId, batchUpdateReq).Do()
	if err != nil {
		log.Fatalf("ドキュメントのテキストを置換できませんでした: %v", err)
	}
	fmt.Println("テキストの置換が完了しました。")

	// エクスポートするファイルの形式
	exportMimeType := "application/pdf"

	// エクスポートのリクエスト作成
	exportRequest := driveSrv.Files.Export(copyDocId, exportMimeType)

	// エクスポート実行
	response, err := exportRequest.Download()
	if err != nil {
		log.Fatalf("ファイルのエクスポートに失敗しました: %v", err)
	}

	// エクスポート結果を保存するファイル名
	outputFileName := fmt.Sprintf("%s_exported_file.pdf", timestamp)

	// 保存先ディレクトリパスを取得
	exportFolderPath := "export"

	// 保存先ディレクトリが存在しない場合は作成する
	if err := os.MkdirAll(exportFolderPath, 0755); err != nil {
		log.Fatalf("保存先ディレクトリの作成に失敗しました: %v", err)
	}

	// ファイルの保存先パスを作成
	outputFilePath := filepath.Join(exportFolderPath, outputFileName)

	// エクスポート結果を保存
	output, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("ファイルの保存に失敗しました: %v", err)
	}
	defer output.Close()

	// エクスポート結果をファイルに書き込む
	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Fatalf("ファイルの書き込みに失敗しました: %v", err)
	}

	fmt.Println("ファイルのエクスポートが完了しました。")
}
