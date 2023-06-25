# 1. Overview
If a template file with embedded dummy data is prepared in advance, you can create user-specific PDF files by replacing the dummy data with appropriate data.

## 1.1. Directory structure
```sh
├── README.md
├── export
├── go.mod
├── go.sum
├── main.go
├── module
│   ├── gdocs.go
│   ├── gdrive.go
│   └── gsheets.go
├── secret.json
├── task
│   ├── gdocs-export.go
│   └── gsheets-export.go
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
- A template file created in Google Sheets.

## 2.1. Set up your environment
### Enable the API
1. Enable the following APi in the Google Cloud console.
  - Google Docs API
  - Google Sheets API
  - Google Drive API.

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
...
...
% 
% gdrive list --service-account secret.json -c . 
Id                                             Name                                   Type   Size     Created
1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q   2023-06-25-07-41-40_Copy-of-Document   doc    4.2 KB   2023-06-25 07:41:40
% 
% gdrive export --force --mime 'application/pdf' --service-account secret.json -c . 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q
Exported '2023-06-25-07-41-40_Copy-of-Document.pdf' with mime type: 'application/pdf'
% gdrive delete --service-account secret.json -c . 1MAFvYn4PAjWXoAfmHGJ4b6NqQeHLtd0LB0YYUPwm17Q
Deleted '2023-06-25-07-41-40_Copy-of-Document'
% 
```