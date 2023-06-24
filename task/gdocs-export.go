package task

import (
	"context"
	"fmt"
	"go-reporting/module"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v2"
)

func GdocsExport() {
	// ========================================
	// 0. 事前準備: Googleドキュメント関連定義
	// ========================================
	// タイムスタンプを取得（現在時刻をJSTに変換）
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp := time.Now().In(jst).Format("2006-01-02-15-04-05")
	// 複製先のGoogleドキュメントのタイトル
	newDocumentTitle := fmt.Sprintf("%s_Copy-of-Document", timestamp)

	// ========================================
	// 0. 事前準備: エクスポートファイル関連定義
	// ========================================
	// エクスポート結果を保存するファイル名
	outputFileName := fmt.Sprintf("%s_exported_file.pdf", timestamp)
	// エクスポートするファイルの形式
	exportMimeType := "application/pdf"
	// 保存先ディレクトリパスを取得
	exportFolderPath := "export"
	// 保存先ディレクトリが存在しない場合は作成する
	if err := os.MkdirAll(exportFolderPath, 0755); err != nil {
		log.Fatalf("保存先ディレクトリの作成に失敗しました: %v", err)
	}
	// ファイルの保存先パスを作成
	outputFilePath := filepath.Join(exportFolderPath, outputFileName)

	// ========================================
	// 0. 事前準備: サービスアカウントのクライアント作成
	// ========================================
	ctx := context.Background()
	// サービスアカウントの秘密鍵を読み込む
	b, err := ioutil.ReadFile("secret.json")
	if err != nil {
		log.Fatalf("秘密鍵ファイルを読み込めませんでした: %v", err)
	}
	// サービスアカウントのクライアントを作成する
	docsSrv, err := module.NewDocsService(ctx, b)
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	driveSrv := module.NewDriveService(ctx, b)
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}

	// ========================================
	// 1. Googleドキュメント複製
	// ========================================
	// 複製リクエストを作成
	copyRequest := &drive.File{
		Title: newDocumentTitle,
	}
	// 環境変数からファイルIDを取得
	sourceDocId := os.Getenv("DOC_ID")
	if sourceDocId == "" {
		// コピー元のドキュメントID
		// https://docs.google.com/document/d/1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg/edit
		sourceDocId = "1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg"
	}
	// Googleドキュメントを複製
	copiedDocument := driveSrv.FileCopy(sourceDocId, copyRequest)
	copyDocId := copiedDocument.Id
	// 複製先のGoogleドキュメントのIDを出力
	fmt.Printf("Googleドキュメントの複製が完了しました。複製先のドキュメントID: %s\n", copyDocId)

	// ========================================
	// 2. Googleドキュメント一覧確認
	// ========================================
	// ファイル一覧を表示する
	driveSrv.FileList()

	// ========================================
	// 3. Googleドキュメントの置換
	// ========================================
	fullName := os.Getenv("FULL_NAME")
	if fullName == "" {
		fullName = "山田 太郎"
	}
	email := os.Getenv("EMAIL")
	if email == "" {
		email = "taro.yamada@test.com"
	}
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
	_, err = docsSrv.Documents.BatchUpdate(copyDocId, batchUpdateReq).Do()
	if err != nil {
		log.Fatalf("ドキュメントのテキストを置換できませんでした: %v", err)
	}
	fmt.Println("テキストの置換が完了しました。")

	// ========================================
	// 4. Googleドキュメントのエクスポート
	// ========================================
	// エクスポート実行
	driveSrv.FileExport(copyDocId, exportMimeType, outputFilePath)
}
