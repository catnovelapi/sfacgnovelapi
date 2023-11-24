package sfclient

import (
	"encoding/json"
	"fmt"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/sfclient/sflogger"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/sfsettings"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"sync"
)

type SFRequest struct {
	m             *sync.Mutex
	debug         bool
	fileLog       *os.File
	builderClient *resty.Client
	Settings      *sfsettings.Settings
}

func NewReqClient() *SFRequest {
	fileInfo, err := os.Stat("sfacgapi.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			// Other error occurred
			log.Println(err)
		}
	} else {
		if fileInfo.Size() > 1024*1024 {
			newName := "sfacgapi" + fileInfo.ModTime().Format("20060102") + ".txt"
			if err = os.Rename("sfacgapi.txt", newName); err != nil {
				log.Println(err)
			}
		}
	}
	// Always open or create "sfacgapi.txt"
	file, err := os.OpenFile("sfacgapi.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open file error : ", err)
	}
	return &SFRequest{
		fileLog:       file,
		debug:         false,
		m:             &sync.Mutex{},
		builderClient: resty.New(),
		Settings:      sfsettings.NewSettings(),
	}
}
func (s *SFRequest) ChangeLogDebug() {
	if s.debug {
		s.debug = false
	} else {
		s.debug = true
	}
}

func (s *SFRequest) SetAuth(username string, password string) {
	s.builderClient.SetBasicAuth(username, password)
}
func (s *SFRequest) SetProxy(proxy string) {
	s.builderClient.SetProxy(proxy)
}
func (s *SFRequest) SetRetryCount(retryCount int) {
	s.builderClient.SetRetryCount(retryCount)
}
func (s *SFRequest) newRequest(params Query, app bool) *resty.Request {
	s.m.Lock()
	defer s.m.Unlock()
	if app {
		s.builderClient.SetBaseURL(s.Settings.GetBaseAPI())
	} else {
		s.builderClient.SetBaseURL(s.Settings.GetWebBaseAPI())
	}
	r := s.builderClient.R().SetHeaders(s.getHeaders())
	if params != nil {
		r.SetQueryParams(params.ToStringMap())
	}
	return r
}
func (s *SFRequest) newLogger(rep *resty.Response) {
	if s.debug {
		if s.fileLog != nil {
			s.m.Lock()
			defer s.m.Unlock()
			if _, err := s.fileLog.WriteString(sflogger.NewSFLogger(rep).CreateLogInfo()); err != nil {
				fmt.Println(err)
			}
		}
	}
}
func (s *SFRequest) Get(endURL string, params Query) (gjson.Result, error) {
	rep, err := s.newRequest(params, true).Get(endURL)
	defer s.newLogger(rep)
	if err != nil {
		return gjson.Result{}, err
	} else if rep.StatusCode() != 200 {
		return gjson.Result{}, fmt.Errorf("status error: %s", rep.Status())
	}
	return gjson.Parse(rep.String()), nil
}

func (s *SFRequest) Post(endURL string, params Query) (*resty.Response, error) {
	var jsonBody string
	if marshal, err := json.Marshal(params); err != nil {
		jsonBody = "{}"
	} else {
		jsonBody = string(marshal)
	}
	rep, err := s.newRequest(nil, true).SetBody(jsonBody).Post(endURL)
	defer s.newLogger(rep)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode() != 200 {
		return nil, fmt.Errorf("status error: %s", rep.Status())
	}
	return rep, nil
}

func (s *SFRequest) PostWeb(endURL string, params Query) (gjson.Result, error) {
	response, err := s.newRequest(params, false).Post(endURL)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Parse(response.String()), nil
}

func (s *SFRequest) GetWeb(endURL string, params Query) (gjson.Result, error) {
	response, err := s.newRequest(params, false).Get(endURL)
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Parse(response.String()), nil
}
func (s *SFRequest) DownloadCover(coverUrl string) ([]byte, error) {
	response, err := s.builderClient.R().SetHeader("User-Agent", s.Settings.GetUserAgent()).Get(coverUrl)
	if err != nil {
		return nil, err
	} else if response.StatusCode() != 200 {
		return nil, fmt.Errorf("status error: %s", response.Status())
	} else {
		return response.Body(), nil
	}
}
