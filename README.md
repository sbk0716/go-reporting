# 1. Overview
- The project allows users to create user-specific PDF files by replacing dummy data embedded in template files created in Google Docs with appropriate data.

## 1.1. Directory structure
```sh
.
├── README.md
├── go.mod
├── go.sum
├── main.go
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
Googleドキュメントの複製が完了しました。複製先のドキュメントID: 1iu7TdR6JmNvLim6T9M3jG_VbxjcPLe0IT5sFrsxKx0U
ファイル一覧:
ファイル名: 2023-06-24-15-36-19_Copy-of-Document (ID: 1iu7TdR6JmNvLim6T9M3jG_VbxjcPLe0IT5sFrsxKx0U)
ファイル名: sample-contract (ID: 1WSzGhnr4rIBVHSTxf1g2bioWarfDtDDhxq1VepMdLwg)
テキストの置換が完了しました。
ファイルのエクスポートが完了しました。
% 
% gdrive list --service-account secret.json -c . 
Id                                             Name                                   Type   Size     Created
1iu7TdR6JmNvLim6T9M3jG_VbxjcPLe0IT5sFrsxKx0U   2023-06-24-15-36-19_Copy-of-Document   doc    4.2 KB   2023-06-24 15:36:22
% 
% gdrive export --force --service-account secret.json -c . 1iu7TdR6JmNvLim6T9M3jG_VbxjcPLe0IT5sFrsxKx0U
Exported '2023-06-24-15-36-19_Copy-of-Document.pdf' with mime type: 'application/pdf'
% gdrive delete --service-account secret.json -c . 1iu7TdR6JmNvLim6T9M3jG_VbxjcPLe0IT5sFrsxKx0U
Deleted '2023-06-24-15-36-19_Copy-of-Document'
% 
```