package task

import (
	"context"
	"encoding/csv"
	"fmt"
	"go-reporting/module"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/drive/v2"
)

func GsheetsBankStatementExport() {
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
	// docsSrv := module.NewDocsService(ctx, b)
	driveSrv := module.NewDriveService(ctx, b)
	sheetsSrv := module.NewSheetsService(ctx, b)

	// ========================================
	// 1. GoogleDrive: ファイル複製
	// ========================================
	fmt.Printf("1. GoogleDrive: ファイル複製\n")
	// 複製リクエストを作成
	copyRequest := &drive.File{
		Title: newFileTitle,
	}
	// スプレッドシートIDとシート名の設定
	sheetName := "main"
	// 環境変数からファイルIDを取得
	spreadsheetId := os.Getenv("SHEET_ID")
	if spreadsheetId == "" {
		// https://docs.google.com/spreadsheets/d/1EobgsZMdxGW5h_3KS_Zjkes3eiUyLOU33rB2_GWi2ss/edit#gid=1792146230
		spreadsheetId = "1EobgsZMdxGW5h_3KS_Zjkes3eiUyLOU33rB2_GWi2ss"
	}
	// GoogleDrive: ファイル複製
	copiedFile := driveSrv.FileCopy(spreadsheetId, copyRequest)
	copyFileId := copiedFile.Id
	fmt.Printf("GoogleDriveにてファイルの複製が完了しました。[ファイルID: %s]\n", copyFileId)

	// ========================================
	// 2. GoogleDrive: ファイル一覧確認
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("2. GoogleDrive: ファイル一覧確認\n")
	// ファイル一覧を表示する
	driveSrv.FileList()

	// ========================================
	// 3. GoogleSheets: データ転記
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("3. GoogleSheets: データ転記\n")

	// CSVファイルを読み込む
	csvFile, err := os.Open("sample.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	// CSVファイルをパースしてデータを取得
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}
	sheetsSrv.TransferDataToSheet(copyFileId, sheetName, records)

	// ========================================
	// 4. GoogleDrive: ファイルエクスポート
	// ========================================
	fmt.Printf("\n")
	fmt.Printf("4. GoogleDrive: ファイルエクスポート\n")
	// エクスポート実行
	driveSrv.FileExport(copyFileId, exportMimeType, outputFilePath)
}
