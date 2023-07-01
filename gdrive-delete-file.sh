#!/bin/bash

# gdrive list コマンドを実行し、結果を変数に格納する
result=$(gdrive list --service-account -c . secret.json --no-header --max 100)
# 結果が空の場合、処理を正常終了する
if [ -z "$result" ]; then
  echo "結果が空です。処理を正常終了します。"
  exit 0
fi

# 結果からIDを抽出してループ処理を行う
while read -r line; do
  # 行からIDを抽出する
  id=$(echo "$line" | awk '{print $1}')

  echo "target file id: $id"
  # id を使用して実行したいコマンドを実行する
  gdrive delete --service-account -c . secret.json "$id"

  # 必要に応じて、他の処理を追加する
done <<< "$result"