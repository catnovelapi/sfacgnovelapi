Here is a detailed GitHub markdown document for the SFAPI code:

# SFAPI

## Overview

The SFAPI provides an API client for interacting with the SFACG website. It allows retrieving book information, searching for novels, getting user account info etc.

## Usage

### Import

```go
import "github.com/catnovelapi/sfapi"
```

### Create Client

```go
client := sfapi.NewSfClient()
```

### Client Options

Set debug mode:

```go 
client.SetDebug()
```

Set cookie:

```go
client.SetCookie(sfCommunity, sessionApp) 
```

Set proxy:

```go
client.SetProxy(proxyURL)
``` 

Set device token:

```go
client.SetDeviceToken(token)
```

Set auth:

```go
client.SetAuth(username, password)
```

Set user agent:

```go
client.SetUserAgent(userAgent)
```

Set retry count:

```go 
client.SetRetryCount(count)
```

Set base URL:

```go
client.SetBaseURL(url)
```

Set device ID:

```go
client.SetDeviceId(deviceID) 
```

Set API version:

```go
client.SetVersion(version)
```

## API Methods

### Get Book Info

```go
client.SF.GetBookInfoApi(bookID)
```

Parameters:

- bookID: book ID

Returns book info in `gjson.Result`

### Get Account Info

```go
client.SF.GetAccountInApi() 
```

Returns user account info in `gjson.Result`

### Get Money Info

```go
client.SF.AccountInMoneyApi()
```

Returns user money info in `gjson.Result`

### Login

```go
client.SF.LoginUsernameApi(username, password)
```

Parameters:

- username: username
- password: password

Returns `*resty.Response`

### Search Novels

```go
client.SF.SearchNovelsResultApi(keyword, page)  
```

Parameters:

- keyword: search keyword
- page: result page number

Returns search result in `gjson.Result`

### Get Chapter Content

```go
client.SF.GetContentInfoApi(chapterID)
```

Parameters:

- chapterID: chapter ID

Returns chapter content in `gjson.Result`

### Get Chapter List

```go 
client.SF.GetChapterDirApi(bookID)
```

Parameters:

- bookID: book ID

Returns chapter list in `gjson.Result`

### Get Bookshelf

```go
client.SF.GetBookShelfApi()  
```

Returns bookshelf info in `gjson.Result`

### Get Book List

```go
client.SF.BookListApi(bookID)
``` 

Parameters:

- bookID: book ID

Returns book list in `gjson.Result`

### Update Books List

```go
client.SF.UpdateBooksList(page)
```

Parameters:

- page: page number

Returns updated book list in `gjson.Result`

### Get Ad Works

```go
client.SF.AdpworksApi(bookID)
```

Parameters:

- bookID: book ID

Returns ad works in `gjson.Result`

### Get Position

```go
client.SF.GetPositionApi()
``` 

Returns position info in `gjson.Result`

### Get Special Push

```go
client.SF.GetSpecialPushApi() 
```

Returns special push info in `gjson.Result`


### Get Welfare Config

```go
client.SF.GetWelfareCfgApi()
```

Returns welfare config in `gjson.Result`

### Get Static Resources

```go
client.SF.GetStaticsResourceApi()
``` 

Returns static resources in `gjson.Result`

### Get User Welfare Store

```go
client.SF.GetUserWelfareStoreitemsLatestApi()
```

Returns user welfare store info in `gjson.Result`

### Get Essay Novels

Get short novels:

```go
client.SF.EssayShortNovelApi(page)  
```

Get novellas:

```go
client.SF.EssayNovellaApi(page)
```

Get long novels:

```go 
client.SF.EssayLongNovelApi(page)
```

Parameters:

- page: page number

Returns essay novels in `gjson.Result`

### Get System Recommend

```go
client.SF.SystemRecommendApi()
```

Returns system recommended books in `gjson.Result`

### Get Act Config

```go 
client.SF.GetActConfigApi()
```

Returns act config in `gjson.Result`

### Post Conversions

```go
client.SF.PostConversionsApi() 
```

Returns conversions result in `gjson.Result`

### Get Version Info

```go
client.SF.VersionInformation()
```

Returns version info in `gjson.Result`

### Pre Order

```go
client.SF.PreOrderApi()
```

Returns `*resty.Response`

### Post Special Push

```go
client.SF.PostSpecialPushApi()
``` 

Returns `*resty.Response`

## APP Methods

### Get Chapter List

```go
client.APP.GetChapterDirList(bookID)
```

Parameters:

- bookID: book ID

Returns chapter list in `[]gjson.Result`

### Get Chapter Content

```go
client.APP.GetChapterContentText(chapterID)
```

Parameters:

- chapterID: chapter ID

Returns chapter content text as `string`

### Get Cookie

```go
client.APP.GetCookie(username, password)  
```

Parameters:

- username: username
- password: password

Returns cookie `string`

### Search Novels

```go
client.APP.GetSearchNovelsResult(keyword, page)
```

Parameters:

- keyword: search keyword
- page: page number

Returns search result novels in `[]gjson.Result`