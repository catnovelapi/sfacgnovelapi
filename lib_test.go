package sfacgnovelapi

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/joho/godotenv"
)

var client *SFClient

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
	client = NewSfClient().R().SetDebug()

}
func TestNewSfClient(t *testing.T) {
	result, err := client.App.GetBookListBySearchKeyword(os.Getenv("BOOK_NAME"), 0)
	if err != nil {
		return
	}
	for index, i := range result {
		fmt.Println(index, "  ", i.Get("novelName").String())
	}
}

func TestNewSfBookInfo(t *testing.T) {
	result, err := client.App.GetBookInfo(os.Getenv("BOOK_ID"))
	if err != nil {
		return
	}
	t.Log(result)
}
func TestNewSfBookChapter(t *testing.T) {
	result, err := client.App.GetChapterContentText(os.Getenv("CHAPTER_ID"))
	if err != nil {
		return
	}
	t.Log(result)
}
func TestNewSfBookInfoTag(t *testing.T) {
	bookList, err := client.App.GetBookListByTag(os.Getenv("TAG_NAME"), 0)
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, book := range bookList {
		wg.Add(1)
		go func(book gjson.Result, wg *sync.WaitGroup) {
			defer wg.Done()
			_, err = client.App.GetChapterDirList(book.Get("novelId").String(), true, nil)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(book.Get("novelName").String(), "加载完成")
			}
		}(book, &wg)
	}
	wg.Wait()

}
