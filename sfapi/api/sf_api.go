package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/sfclient"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/semaphore"
	"sync"
)

type API struct {
	sync.Mutex
	sem *semaphore.Weighted
	Req *sfclient.SFRequest
}

func NewApi() *API {
	return &API{
		sem: semaphore.NewWeighted(32),
		Req: sfclient.NewReqClient(),
	}
}

const defaultBookInfoExpand = "chapterCount,bigBgBanner,bigNovelCover,typeName,intro,fav,ticket,pointCount,tags,sysTags,signlevel,discount,discountExpireDate,totalNeedFireMoney,rankinglist,originTotalNeedFireMoney,firstchapter,latestchapter,latestcommentdate,essaytag,auditCover,preOrderInfo,customTag,topic,unauditedCustomtag,homeFlag,isbranch,essayawards"

// SetNewSem sets the maximum number of concurrent requests
func (sfacg *API) SetNewSem(maxConcurrentRequests int) {
	sfacg.Lock()
	defer sfacg.Unlock()
	sfacg.sem = semaphore.NewWeighted(int64(maxConcurrentRequests))
}

// GetBookInfoApi get book information by bookId
func (sfacg *API) GetBookInfoApi(ctx context.Context, bookId string) (gjson.Result, error) {
	// Acquire a token, limiting the number of concurrent requests
	if err := sfacg.sem.Acquire(ctx, 1); err != nil {
		return gjson.Result{}, err
	}
	defer sfacg.sem.Release(1)
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v", bookId), sfclient.Query{"expand": defaultBookInfoExpand})
}

func (sfacg *API) GetBookCommentBarrageNewVersionApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/cmts/barragenewversion", bookId), nil)
}
func (sfacg *API) GetBookCommentApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/Cmts", bookId), sfclient.Query{"page": 0, "size": 2, "type": "stickandclear", "sort": "timeline", "replyUserId": 0})
}

func (sfacg *API) GetBookCommentsApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/lcmts", bookId), sfclient.Query{"page": 0, "size": 1, "sort": "addtime", "charlen": 140})
}

func (sfacg *API) GetBookViewsApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/user/NovelViews/%v", bookId), nil)
}
func (sfacg *API) GetBookBonusRankApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/bonus/rank", bookId), sfclient.Query{"numMax": 50, "dateRange": 1})
}

func (sfacg *API) GetBookTicketRankApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/ticket/rank", bookId), sfclient.Query{"numMax": 50, "dateRange": 1})
}
func (sfacg *API) GetBookSignlevelApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/adpworks/novelId/%v", bookId), sfclient.Query{"expand": "signlevel"})
}

func (sfacg *API) BookListApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novel/%v/bookList", bookId), sfclient.Query{"size": 3, "page": 0})
}

func (sfacg *API) BookFansList(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/fans", bookId), nil)
}

// GetAccountInApi get account information
func (sfacg *API) GetAccountInApi() (gjson.Result, error) {
	return sfacg.Req.Get("/user", sfclient.Query{"expand": "vipInfo,welfareMoney,newVip"})
}

// GetAccountInMoneyApi get account money information
func (sfacg *API) GetAccountInMoneyApi() (gjson.Result, error) {
	return sfacg.Req.Get("/user/money", nil)
}

// LoginUsernameApi login by username and password
func (sfacg *API) LoginUsernameApi(username, password string) (*resty.Response, error) {
	return sfacg.Req.Post("/sessions", sfclient.Query{"username": username, "password": password})
}

func (sfacg *API) SearchNovelsResultApi(keyword string, page string) (gjson.Result, error) {
	query := sfclient.Query{
		"q":          keyword,
		"page":       page,
		"size":       "12",
		"sort":       "hot",
		"systagids":  "",
		"isFinish":   "-1",
		"updateDays": "-1",
		"expand":     defaultBookInfoExpand}
	return sfacg.Req.Get("/search/novels/result/new", query)
}

func (sfacg *API) GetContentInfoApi(ctx context.Context, chapterId string) (gjson.Result, error) {
	// Acquire a token, limiting the number of concurrent requests
	if err := sfacg.sem.Acquire(ctx, 1); err != nil {
		return gjson.Result{}, err
	}
	defer sfacg.sem.Release(1)

	query := sfclient.Query{"expand": "content,needFireMoney,originNeedFireMoney,tsukkomi,chatlines,isbranch,isContentEncrypted,authorTalk", "autoOrder": "false"}
	if response, err := sfacg.Req.Get(fmt.Sprintf("/Chaps/%v", chapterId), query); err != nil {
		return gjson.Result{}, err
	} else if response.Get("data.expand.content").String() == "" {
		return gjson.Result{}, fmt.Errorf("chapterId:%v,获取章节内容失败", chapterId)
	} else {
		return response.Get("data"), nil
	}
}

func (sfacg *API) GetActpushesApi(bookId any) (gjson.Result, error) {
	return sfacg.Req.Get(fmt.Sprintf("/novels/%v/actpushes", bookId), sfclient.Query{"filter": "android", "pageType": 1})
}
func (sfacg *API) GetSpecialpushsNewNovelApi(page any) (gjson.Result, error) {
	return sfacg.Req.Get("/novels/specialpushs", sfclient.Query{"pushNames": "newpush", "page": page, "size": 8, "expand": defaultBookInfoExpand})
}

func (sfacg *API) GetSpecialpushsHotNovelApi(page any) (gjson.Result, error) {
	return sfacg.Req.Get("/novels/specialpushs", sfclient.Query{"pushNames": "hotpush", "page": page, "size": 8, "expand": defaultBookInfoExpand})
}
func (sfacg *API) GetChapterDirApi(ctx context.Context, bookId string) (gjson.Result, error) {
	// Acquire a token, limiting the number of concurrent requests
	if err := sfacg.sem.Acquire(ctx, 1); err != nil {
		return gjson.Result{}, err
	}
	defer sfacg.sem.Release(1)

	result, err := sfacg.Req.Get(fmt.Sprintf("/novels/%s/dirs", bookId), sfclient.Query{"expand": "originNeedFireMoney"})
	if err != nil {
		return gjson.Result{}, err
	}
	return result.Get("data"), nil
}
func (sfacg *API) GetBookShelfApi() (gjson.Result, error) {
	return sfacg.Req.Get("/user/Pockets", sfclient.Query{"expand": "novels"})
}

func (sfacg *API) GetBookListByTagApi(tagId, page any) (gjson.Result, error) {
	query := sfclient.Query{"page": fmt.Sprintf("%v", page), "size": "20", "expand": defaultBookInfoExpand, "sort": "viewtimes", "filter": "all"}
	return sfacg.Req.Get(fmt.Sprintf("/novels/0/sysTags/%v/novels", tagId), query)
}
func (sfacg *API) GetTagNameApi() (gjson.Result, error) {
	result, err := sfacg.Req.Get("/novels/0/sysTags", sfclient.Query{"categoryId": 0})
	if err != nil {
		return gjson.Result{}, err
	}
	return result.Get("data"), nil
}
func (sfacg *API) GetRanksWeekNovelsApi(rankType string, page any) (gjson.Result, error) {
	return sfacg.Req.Get("/ranks/week/novels", sfclient.Query{"page": page, "size": 50, "rtype": rankType, "ntype": "origin", "expand": defaultBookInfoExpand})
}
func (sfacg *API) GetRanksMonthNovelsApi(rankType string, page any) (gjson.Result, error) {
	return sfacg.Req.Get("/ranks/month/novels", sfclient.Query{"page": page, "size": 50, "rtype": rankType, "ntype": "origin", "expand": defaultBookInfoExpand})
}
func (sfacg *API) GetRanksAllNovelsApi(rankType string, page any) (gjson.Result, error) {
	return sfacg.Req.Get("/ranks/all/novels", sfclient.Query{"page": page, "size": 50, "rtype": rankType, "ntype": "origin", "expand": defaultBookInfoExpand})
}
func (sfacg *API) UpdateBooksList(page string) (gjson.Result, error) {
	return sfacg.Req.Get("/novels", sfclient.Query{"page": page, "size": 50, "filter": "latest-signnvip", "expand": "sysTags,intro"})
}

func (sfacg *API) GetPositionApi() (gjson.Result, error) {
	return sfacg.Req.Get("/position", nil)
}

func (sfacg *API) GetSpecialPushApi() (gjson.Result, error) {
	return sfacg.Req.Get("/specialpush", sfclient.Query{"pushNames": "merchPush", "entityId": "", "entityType": "novel"})
}

func (sfacg *API) GetWelfareCfgApi() (gjson.Result, error) {
	return sfacg.Req.Get("/welfare/cfg", nil)
}

func (sfacg *API) GetStaticsResourceApi() (gjson.Result, error) {
	return sfacg.Req.Get("/StaticsResource", nil)
}

func (sfacg *API) GetUserWelfareStoreitemsLatestApi() (gjson.Result, error) {
	return sfacg.Req.Get("/user/welfare/storeitems/latest", nil)
}

func (sfacg *API) essaySolicitationNovelApi(tagIds, page string) (gjson.Result, error) {
	return sfacg.Req.GetWeb("/api/essay/getNovels", sfclient.Query{"tagIds": tagIds, "page": page, "size": 50})
}

func (sfacg *API) EssayShortNovelApi(page any) (gjson.Result, error) {
	return sfacg.essaySolicitationNovelApi("655", fmt.Sprintf("%v", page))
}

func (sfacg *API) EssayNovellaApi(page any) (gjson.Result, error) {
	return sfacg.essaySolicitationNovelApi("654", fmt.Sprintf("%v", page))
}

func (sfacg *API) EssayLongNovelApi(page int) (gjson.Result, error) {
	return sfacg.essaySolicitationNovelApi("653", fmt.Sprintf("%v", page))
}

func (sfacg *API) SystemRecommendApi() (gjson.Result, error) {
	return sfacg.Req.Get("/novel/systemRecommend", nil)
}

func (sfacg *API) GetActConfigApi() (gjson.Result, error) {
	return sfacg.Req.GetWeb("/api/apptabproject/getActConfig", nil)
}
func (sfacg *API) PostConversionsApi() (gjson.Result, error) {
	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	return sfacg.Req.Get("/androiddeviceinfos/conversion", sfclient.Query{"oaID": hex.EncodeToString(randomBytes)})

}
func (sfacg *API) VersionInformation() (gjson.Result, error) {
	return sfacg.Req.Get("/androidcfg", nil)
}
func (sfacg *API) PreOrderApi() (*resty.Response, error) {
	return sfacg.Req.Post("/preOrder", sfclient.Query{"expand": "intro,typeName,tags,sysTags", "withExpiredPreOrder": "false"})
}
func (sfacg *API) PostSpecialPushApi() (*resty.Response, error) {
	return sfacg.Req.Post("/specialpush", sfclient.Query{"signDate": "merchPush", "entityId": "", "entityType": "novel"})
}
