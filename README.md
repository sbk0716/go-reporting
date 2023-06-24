# 1. Overview
- The project allows users to create user-specific PDF files by replacing dummy data embedded in template files created in Google Docs with appropriate data.

## 1.1. Directory structure
```sh
.
├── README.md
├── export
├── go.mod
├── go.sum
├── main.go
├── module
│   ├── gdocs.go
│   └── gdrive.go
├── secret.json
├── task
│   └── gdocs-export.go
└── template
    └── sample-contract.pdf
```


# 2. Usage
## 2.0. Prerequisites
- Latest version of Go.
- Latest version of Git.
- A Google Cloud project.
- A Google Account.
- A template file created in Google Docs.

## 2.1. Set up your environment
### Enable the API
1. In the Google Cloud console, enable the Google Docs API and Google Drive API.

### Configure the OAuth consent screen

1. Refer to the following URLs to set up.
  - https://developers.google.com/docs/api/quickstart/go#configure_the_oauth_consent_screen

### Create the service account
1. In the Google Cloud console, go to Menu menu > APIs & Services > Credentials.
2. Click Create Credentials > Service account.
3. Enter the service account details to create the service account.
4. Create the private key for the service account.
5. Save the private key as secret.json, and move the file to your working directory.


## 2.2. Run app

```sh
% go get .
% go run main.go
```

# 3. Check the operation

```sh
% gdrive list --service-account secret.json -c .
Id   Name   Type   Size   Created
% 
% go run main.go
1. Googleドキュメントの複製
Googleドキュメントの複製が完了しました。複製先のドキュメントID: 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q

2. Googleドキュメント一覧確認
### [ファイル一覧] ###
ファイル名: 2023-06-25-07-41-40_Copy-of-Document (ID: 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q)
ファイル名: sample-contract (ID: 1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg)

3. Googleドキュメントの置換
ReplaceAllText find: "${fullName}"
ReplaceAllText replace: "山田 太郎"
ReplaceAllText find: "${email}"
ReplaceAllText replace: "taro.yamada@test.com"
テキストの置換が完了しました。

4. Googleドキュメントのエクスポート
ファイルのエクスポートが完了しました。
% 
% gdrive list --service-account secret.json -c . 
Id                                             Name                                   Type   Size     Created
1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q   2023-06-25-07-41-40_Copy-of-Document   doc    4.2 KB   2023-06-25 07:41:40
% 
% gdrive export --force --service-account secret.json -c . 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q
Exported '2023-06-25-07-41-40_Copy-of-Document.pdf' with mime type: 'application/pdf'
% gdrive delete --service-account secret.json -c . 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q
Deleted '2023-06-25-07-41-40_Copy-of-Document'
% 
```