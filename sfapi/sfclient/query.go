package sfclient

import "fmt"

type Query map[string]any

func Q() Query {
	return make(Query)
}

func (q Query) Encode() string {
	var s string
	for k, v := range q {
		s += fmt.Sprintf("%v=%v&", k, v)
	}
	return s[:len(s)-1]
}
func (q Query) Add(key string, value any) {
	q[key] = value
}
func (q Query) AddMap(m map[string]any) {
	for k, v := range m {
		q[k] = v
	}
}
func (q Query) Get(key string) any {
	return q[key]
}
func (q Query) Del(key string) {
	if _, ok := q[key]; ok {
		delete(q, key)
	}
}

func (q Query) SetBookInfoQuery() Query {
	const defaultBookInfoExpand = "chapterCount,bigBgBanner,bigNovelCover,typeName,intro,fav,ticket,pointCount,tags,sysTags,signlevel,discount,discountExpireDate,totalNeedFireMoney,rankinglist,originTotalNeedFireMoney,firstchapter,latestchapter,latestcommentdate,essaytag,auditCover,preOrderInfo,customTag,topic,unauditedCustomtag,homeFlag,isbranch,essayawards"
	if _, ok := q["expand"]; !ok {
		q["expand"] = defaultBookInfoExpand
	} else {
		q["expand"] = fmt.Sprintf("%v,%v", q["expand"], defaultBookInfoExpand)
	}
	return q
}

func (q Query) SetPageQuery(page any) Query {
	q["page"] = fmt.Sprintf("%v", page)
	q["size"] = 50
	return q
}

func (q Query) ToStringMap() map[string]string {
	m := make(map[string]string)
	for k, v := range q {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}
