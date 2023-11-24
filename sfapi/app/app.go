package app

import (
	"context"
	"fmt"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/api"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
)

type APP struct {
	API *api.API
}

func NewApp(api *api.API) *APP {
	return &APP{API: api}
}

// GetChapterDirList get chapter list bookId:book id, vipVerification:whether vip verification is required, verificationFunc:verification function, return true to skip
func (sfacg *APP) GetChapterDirList(bookId any, vipVerification bool, verificationFunc func(gjson.Result) bool) ([]gjson.Result, error) {
	response, err := sfacg.API.GetChapterDirApi(context.TODO(), fmt.Sprintf("%v", bookId))
	if err != nil {
		return nil, err
	}
	var chapterList []gjson.Result
	for _, i := range response.Get("volumeList").Array() {
		for _, j := range i.Get("chapterList").Array() {
			if verificationFunc != nil {
				if verificationFunc(j) {
					continue
				}
			}
			if vipVerification {
				if j.Get("originNeedFireMoney").Int() != 0 {
					continue
				}
			}
			chapterList = append(chapterList, j)
		}
	}
	if len(chapterList) == 0 {
		return nil, fmt.Errorf("bookId:%v,章节列表为空", bookId)
	}
	return chapterList, nil
}

func (sfacg *APP) GetChapterContentText(chapterId any) (string, error) {
	response, err := sfacg.API.GetContentInfoApi(context.TODO(), fmt.Sprintf("%v", chapterId))
	if err != nil {
		return "", err
	}
	var content = response.Get("ntitle").String() + "\n"
	for _, i := range strings.Split(response.Get("expand.content").String(), "\n") {
		if j := strings.ReplaceAll(strings.TrimSpace(i), "　", ""); j != "" {
			content += fmt.Sprintf("　　%v\n", j)
		}
	}
	return content, nil
}

func addFunc(bookId any, fs ...func(bookId any) (gjson.Result, error)) {
	wg := sync.WaitGroup{}
	for _, i := range fs {
		wg.Add(1)
		go func(i func(bookId any) (gjson.Result, error)) {
			defer wg.Done()
			_, err := i(bookId)
			if err != nil {
				fmt.Println(err)
			}
		}(i)
	}
	wg.Wait()
}
func (sfacg *APP) GetBookInfo(bookId any) (gjson.Result, error) {
	bookInfoApi, err := sfacg.API.GetBookInfoApi(context.TODO(), fmt.Sprintf("%v", bookId))
	if err != nil {
		return gjson.Result{}, err
	} else {
		addFunc(
			bookId,
			sfacg.API.GetActpushesApi,
			sfacg.API.GetBookCommentApi,
			sfacg.API.GetBookCommentBarrageNewVersionApi,
			sfacg.API.GetBookCommentsApi,
			sfacg.API.BookListApi,
			sfacg.API.BookFansList,
			sfacg.API.GetBookTicketRankApi,
			sfacg.API.GetBookViewsApi,
			sfacg.API.GetBookBonusRankApi,
		)
		return bookInfoApi, nil
	}
}
func (sfacg *APP) GetCookie(username, password string) (string, error) {
	sessions, err := sfacg.API.LoginUsernameApi(username, password)
	if err != nil {
		return "", err
	}
	var cookie string
	for _, i := range sessions.Cookies() {
		cookie += i.Name + "=" + i.Value + ";"
	}
	if cookie == "" {
		return "", fmt.Errorf("获取cookie失败, cookie为空")
	}
	// set cookie to sfacg.API.Req.builderClient
	sfacg.API.Req.Settings.SetCookie(sessions.Cookies())
	return cookie, nil

}
func (sfacg *APP) GetBookListBySearchKeyword(keyword string, page any) ([]gjson.Result, error) {
	resultApi, err := sfacg.API.SearchNovelsResultApi(keyword, fmt.Sprintf("%v", page))
	if err != nil {
		return nil, err
	}
	var novels []gjson.Result
	for _, i := range resultApi.Get("data.novels").Array() {
		novels = append(novels, i)
	}
	if len(novels) == 0 {
		return nil, fmt.Errorf("keyword:%v,获取搜索结果失败,结果为空", keyword)
	}
	return novels, nil
}

func (sfacg *APP) GetBookListByTag(tagName string, page any) ([]gjson.Result, error) {
	tagNameList, err := sfacg.API.GetTagNameApi()
	if err != nil {
		return nil, err
	}
	for _, i := range tagNameList.Array() {
		if i.Get("tagName").String() == tagName {
			result, ok := sfacg.API.GetBookListByTagApi(i.Get("sysTagId").Int(), page)
			if ok != nil {
				return nil, ok
			}
			var novels []gjson.Result
			for _, novel := range result.Get("data").Array() {
				novels = append(novels, novel)
			}
			if len(novels) == 0 {
				return nil, fmt.Errorf("tagName:%v,获取标签书籍列表失败,结果为空", tagName)
			}
			return novels, nil
		}
	}
	return nil, fmt.Errorf("tagName:%v,获取标签书籍列表失败,标签不存在", tagName)
}
