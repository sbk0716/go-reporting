package module

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type DriveService struct {
	service *drive.Service
}

// サービスアカウントのクライアントを作成する
func NewDriveService(ctx context.Context, b []byte) *DriveService {
	// サービスアカウントのクライアントを作成する
	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	return &DriveService{
		service: srv,
	}
}

// ファイルを複製する
func (d *DriveService) FileCopy(fileId string, copyRequest *drive.File) *drive.File {
	copiedFile, err := d.service.Files.Copy(fileId, copyRequest).Do()
	if err != nil {
		log.Fatalf("ファイルの複製に失敗しました: %v", err)
	}
	return copiedFile
}

// GoogleDriveのファイル一覧を表示する
func (d *DriveService) FileList() *drive.FileList {
	files, err := d.service.Files.List().Fields("*").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	// GoogleDriveのファイル一覧を表示する
	fmt.Println("### [GoogleDrive: ファイル一覧] ###")
	if len(files.Items) > 0 {
		for _, file := range files.Items {
			fmt.Printf("ファイル名: %s (ID: %s)\n", file.Title, file.Id)
		}
	} else {
		fmt.Println("ファイルが見つかりませんでした。")
	}
	return files
}

// GoogleDriveのファイルをPDFとしてエクスポートする
func (d *DriveService) FileExport(fileId string, mimeType string, outputFilePath string) {
	// エクスポートのリクエスト作成
	exportRequest := d.service.Files.Export(fileId, mimeType)
	// エクスポート実行
	response, err := exportRequest.Download()
	if err != nil {
		log.Fatalf("ファイルのエクスポートに失敗しました: %v", err)
	}
	// エクスポート結果保存用ファイルの作成
	output, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("エクスポート結果保存用ファイルの作成に失敗しました: %v", err)
	}
	defer output.Close()

	// エクスポート結果をファイルに書き込む
	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Fatalf("エクスポート結果のファイルへの書き込みに失敗しました: %v", err)
	}
	fmt.Println("ファイルのエクスポートが完了しました。")
}
