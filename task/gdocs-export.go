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

	"google.golang.org/api/drive/v2"
)

func GdocsExport() {
	// ========================================
	// 0. 事前準備: Googleドキュメント関連定義
	// ========================================
	// タイムスタンプを取得（現在時刻をJSTに変換）
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp := time.Now().In(jst).Format("2006-01-02-15-04-05")
	// 複製先のファイルのタイトル
	newFileTitle := fmt.Sprintf("%s_Copy", timestamp)

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
	docsSrv := module.NewDocsService(ctx, b)
	driveSrv := module.NewDriveService(ctx, b)

	// ========================================
	// 1. GoogleDrive: ファイル複製
	// ========================================
	fmt.Printf("1. GoogleDrive: ファイル複製\n")
	// 複製リクエストを作成
	copyRequest := &drive.File{
		Title: newFileTitle,
	}
	// 環境変数からファイルIDを取得
	docId := os.Getenv("DOC_ID")
	if docId == "" {
		// コピー元のドキュメントID
		// https://docs.google.com/document/d/1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg/edit
		docId = "1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg"
	}
	// GoogleDrive: ファイル複製
	copiedFile := driveSrv.FileCopy(docId, copyRequest)
	copyFileId := copiedFile.Id
	fmt.Printf("GoogleDriveにてファイルの複製が完了しました。[ファイルID: %s]\n", copyFileId)

	// ========================================
	// 2. GoogleDriveのファイル一覧確認
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("2. GoogleDriveのファイル一覧確認\n")
	// ファイル一覧を表示する
	driveSrv.FileList()

	// ========================================
	// 3. GoogleSheets: テキスト置換
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("3. GoogleSheets: テキスト置換\n")
	fullName := os.Getenv("FULL_NAME")
	if fullName == "" {
		fullName = "山田 太郎"
	}
	email := os.Getenv("EMAIL")
	if email == "" {
		email = "taro.yamada@test.com"
	}
	// 置換対象の文字列と置換後の文字列のマップ
	replacements := map[string]string{
		"${fullName}": fullName,
		"${email}":    email,
	}
	docsSrv.ReplaceAllText(docId, replacements)

	// ========================================
	// 4. GoogleDrive: ファイルエクスポート
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("4. ファイルエクスポート\n")
	// エクスポート実行
	driveSrv.FileExport(docId, exportMimeType, outputFilePath)
}
