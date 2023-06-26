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

// 対象のスプレッドシートのシートIDを取得する
func getSheetID(sheet *sheets.Spreadsheet, sheetName string) int64 {
	for _, s := range sheet.Sheets {
		if s.Properties.Title == sheetName {
			return s.Properties.SheetId
		}
	}
	return 0
}

// 対象のスプレッドシートのテキストを置換する
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

// 対象のスプレッドシートにデータを転記する
func (s *SheetsService) TransferDataToSheet(spreadsheetID string, sheetName string, data [][]string) error {
	// 転記するデータを作成
	values := make([][]interface{}, len(data))
	for i, record := range data {
		row := make([]interface{}, len(record))
		for j, value := range record {
			row[j] = value
		}
		values[i] = row
	}

	// 対象のスプレッドシートの情報を取得する
	sheet, err := s.service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("スプレッドシートの情報の取得が失敗しました: %v", err)
		return err
	}

	// シートIDを取得
	sheetID := getSheetID(sheet, sheetName)

	// 転記するデータの範囲を指定
	startRowIndex := int64(9) // 10行目からデータを転記する
	endRowIndex := startRowIndex + int64(len(data))
	startColumnIndex := int64(0) // A列から転記する
	endColumnIndex := startColumnIndex + int64(len(data[0]))

	// 転記するデータをシートに書き込む
	_, err = s.service.Spreadsheets.Values.Update(spreadsheetID, fmt.Sprintf("%s!A%d", sheetName, startRowIndex), &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("スプレッドシートへのデータの転記が失敗しました: %v", err)
		return err
	}
	fmt.Println("スプレッドシートへのデータの転記が完了しました。")

	// 書式設定の範囲を指定
	formattingRange := &sheets.GridRange{
		SheetId:          sheetID,
		StartRowIndex:    startRowIndex,
		EndRowIndex:      endRowIndex,
		StartColumnIndex: startColumnIndex,
		EndColumnIndex:   endColumnIndex,
	}

	// 書式設定をコピーして貼り付けるリクエストを作成
	formattingRequests := []*sheets.Request{}
	formattingRequests = append(formattingRequests, &sheets.Request{
		CopyPaste: &sheets.CopyPasteRequest{
			Source: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    startRowIndex,
				EndRowIndex:      10, // 10行目の書式をコピーする
				StartColumnIndex: startColumnIndex,
				EndColumnIndex:   endColumnIndex,
			},
			Destination: formattingRange,
			PasteType:   "PASTE_FORMAT",
		},
	})

	// バッチ更新リクエストを作成
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: formattingRequests,
	}

	// 書式設定を適用する
	_, err = s.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Do()
	if err != nil {
		log.Fatalf("スプレッドシートの書式設定が失敗しました: %v", err)
		return err
	}
	fmt.Println("スプレッドシートの書式設定が完了しました。")
	return nil
}

func (s *SheetsService) SetPrintRange(spreadsheetID, sheetName string, data [][]string) error {
	// 対象のスプレッドシートの情報を取得する
	sheet, err := s.service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("スプレッドシートの情報の取得が失敗しました: %v", err)
		return err
	}

	// シートIDを取得
	sheetID := getSheetID(sheet, sheetName)
	// 印刷するデータの範囲を指定
	startRowIndex := int64(0) // 1行目からデータを転記する
	endRowIndex := startRowIndex + int64(len(data))
	startColumnIndex := int64(0) // A列から転記する
	endColumnIndex := startColumnIndex + int64(len(data[0]))

	request := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						SheetId:    sheetID,
						Dimension:  "ROWS",
						StartIndex: startRowIndex,
						EndIndex:   endRowIndex,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 20,
					},
					Fields: "pixelSize",
				},
			},
			{
				UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
					Range: &sheets.DimensionRange{
						SheetId:    sheetID,
						Dimension:  "COLUMNS",
						StartIndex: startColumnIndex,
						EndIndex:   endColumnIndex,
					},
					Properties: &sheets.DimensionProperties{
						PixelSize: 200,
					},
					Fields: "pixelSize",
				},
			},
		},
	}

	// 印刷設定を変更する
	_, err = s.service.Spreadsheets.BatchUpdate(spreadsheetID, request).Do()
	if err != nil {
		log.Fatalf("スプレッドシートの印刷範囲の設定が失敗しました: %v", err)
		return err
	}
	fmt.Println("スプレッドシートの印刷範囲の設定が完了しました。")

	return nil
}

func (s *SheetsService) SetScaleToFitWidth(spreadsheetID, sheetName string, data [][]string) error {
	// 対象のスプレッドシートの情報を取得する
	sheet, err := s.service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("スプレッドシートの情報の取得が失敗しました: %v", err)
		return err
	}

	// シートIDを取得
	sheetID := getSheetID(sheet, sheetName)
	// startRowIndex := int64(0) // 1行目からデータを転記する
	// endRowIndex := startRowIndex + int64(len(data))
	startColumnIndex := int64(0) // A列から転記する
	endColumnIndex := startColumnIndex + int64(len(data[0]))

	// スケールの設定を行うリクエストを作成
	requests := []*sheets.Request{
		{
			AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
				Dimensions: &sheets.DimensionRange{
					SheetId:    sheetID,
					Dimension:  "COLUMNS",
					StartIndex: startColumnIndex,
					EndIndex:   endColumnIndex,
				},
			},
		},
	}

	// リクエストをバッチ更新として送信
	_, err = s.service.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Fatalf("スプレッドシートのスケールの設定が失敗しました: %v", err)
		return err
	}
	fmt.Println("スプレッドシートのスケールの設定が完了しました。")

	return nil
}
