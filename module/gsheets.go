package module

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	service *sheets.Service
}

// サービスアカウントのクライアントを作成する
func NewSheetsService(ctx context.Context, b []byte) *SheetsService {
	// サービスアカウントのクライアントを作成する
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("サービスアカウントのクライアントを作成できませんでした: %v", err)
	}
	return &SheetsService{
		service: srv,
	}
}

func (s *SheetsService) ReplaceAllText(spreadsheetID, sheetName string, replacements map[string]string) {
	// 置換対象のセル範囲を指定
	rangeValue := fmt.Sprintf("%s!A1:ZZ", sheetName)
	writeRange := fmt.Sprintf("%s!A1:ZZ", sheetName)

	// スプレッドシートからデータを取得
	resp, err := s.service.Spreadsheets.Values.Get(spreadsheetID, rangeValue).Do()
	if err != nil {
		log.Fatalf("スプレッドシートからデータを取得できませんでした: %v", err)
	}
	// データの置換
	for _, row := range resp.Values {
		for cellIndex, cellValue := range row {
			cellValueStr, ok := cellValue.(string)
			if !ok {
				continue
			}

			// 置換対象の文字列が存在する場合は置換する
			for searchStr, replaceStr := range replacements {
				cellValueStr = strings.ReplaceAll(cellValueStr, searchStr, replaceStr)
			}

			row[cellIndex] = cellValueStr
		}
	}
	// 置換結果をスプレッドシートに書き込む
	valueRange := &sheets.ValueRange{
		Values: resp.Values,
	}
	_, err = s.service.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("スプレッドシートにデータを書き込めませんでした: %v", err)
	}

	fmt.Println("スプレッドシートのテキストの置換が完了しました。")
}
