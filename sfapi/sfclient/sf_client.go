package sfclient

import (
	"encoding/json"
	"fmt"
	"github.com/catnovelapi/builder"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/sfsettings"
	"github.com/tidwall/gjson"
	"sync"
)

type SFRequest struct {
	m             *sync.Mutex
	builderClient *builder.Client
	Settings      *sfsettings.Settings
}

func NewReqClient() *SFRequest {
	return &SFRequest{
		m:             &sync.Mutex{},
		builderClient: builder.NewClient(),
		Settings:      sfsettings.NewSettings(),
	}
}
func (s *SFRequest) LogDebug() {
	s.builderClient.SetDebug()
	s.builderClient.SetDebugFile("sfacgapi.txt")
}

func (s *SFRequest) SetAuth(username string, password string) {
	s.builderClient.SetBasicAuth(username, password)
}
func (s *SFRequest) SetProxy(proxy string) {
	s.builderClient.SetProxy(proxy)
}
func (s *SFRequest) SetRetryCount(retryCount int) {
	s.builderClient.SetRetryNumber(retryCount)
}
func (s *SFRequest) newRequest(params Query, app bool) *builder.Request {
	s.m.Lock()
	defer s.m.Unlock()
	if app {
		s.builderClient.SetBaseURL(s.Settings.GetBaseAPI())
	} else {
		s.builderClient.SetBaseURL(s.Settings.GetWebBaseAPI())
	}
	if s.builderClient.GetClientCookie() == "" {
		if s.Settings.GetCookie() != "" {
			s.builderClient.SetCookie(s.Settings.GetCookie())
		}
	}
	r := s.builderClient.R().SetHeaders(s.getHeaders())
	if params != nil {
		r.SetQueryParams(params)
	}
	return r
}
func (s *SFRequest) Get(endURL string, params Query) (gjson.Result, error) {
	rep, err := s.newRequest(params, true).Get(endURL)
	if err != nil {
		return gjson.Result{}, err
	} else if !rep.IsStatusOk() {
		return gjson.Result{}, fmt.Errorf("status error: %s", rep.GetStatus())
	}
	return gjson.Parse(rep.String()), nil
}

func (s *SFRequest) Post(endURL string, params Query) (*builder.Response, error) {
	var jsonBody string
	if marshal, err := json.Marshal(params); err != nil {
		jsonBody = "{}"
	} else {
		jsonBody = string(marshal)
	}
	rep, err := s.newRequest(nil, true).SetBody(jsonBody).Post(endURL)
	if err != nil {
		return nil, err
	}
	if !rep.IsStatusOk() {
		return nil, fmt.Errorf("status error: %s", rep.GetStatus())
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
	} else if !response.IsStatusOk() {
		return nil, fmt.Errorf("status error: %s", response.GetStatus())
	} else {
		return response.GetByte(), nil
	}
}
