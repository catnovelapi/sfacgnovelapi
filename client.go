package sfacgnovelapi

import (
	"github.com/catnovelapi/sfacgnovelapi/sfapi/api"
	"github.com/catnovelapi/sfacgnovelapi/sfapi/app"
	"net/http"
)

type SFClient struct {
	Api *api.API
	App *app.APP
}

func NewSfClient() *SFClient {
	sfClient := &SFClient{Api: api.NewApi()}
	sfClient.App = app.NewApp(sfClient.Api)
	return sfClient
}
func (sfClient *SFClient) R() *SFClient {
	return sfClient.SetRetryCount(5).SetAuth("androiduser", "1a#$51-yt69;*Acv@qxq")
}

func (sfClient *SFClient) SetDebug() *SFClient {
	sfClient.Api.Req.LogDebug()

	return sfClient
}

func (sfClient *SFClient) SetCookie(sfCommunity string, sessionApp string) *SFClient {
	sfClient.Api.Req.Settings.SetCookie([]*http.Cookie{{Name: ".SFCommunity", Value: sfCommunity}, {Name: "session_APP", Value: sessionApp}})
	return sfClient
}
func (sfClient *SFClient) SetProxy(proxy string) *SFClient {
	sfClient.Api.Req.SetProxy(proxy)
	return sfClient
}
func (sfClient *SFClient) SetNewSem(maxConcurrentRequests int) *SFClient {
	sfClient.Api.SetNewSem(maxConcurrentRequests)
	return sfClient
}
func (sfClient *SFClient) SetAndroidKey(androidKey string) *SFClient {
	sfClient.Api.Req.Settings.SetAndroidKey(androidKey)
	return sfClient
}
func (sfClient *SFClient) SetAuth(username string, password string) *SFClient {
	sfClient.Api.Req.SetAuth(username, password)
	return sfClient
}
func (sfClient *SFClient) SetUserAgent(userAgent string) *SFClient {
	sfClient.Api.Req.Settings.SetUserAgent(userAgent)
	return sfClient
}
func (sfClient *SFClient) SetRetryCount(retryCount int) *SFClient {
	sfClient.Api.Req.SetRetryCount(retryCount)
	return sfClient
}
func (sfClient *SFClient) SetBaseURL(baseURL string) *SFClient {
	sfClient.Api.Req.Settings.SetBaseURL(baseURL)
	return sfClient
}
func (sfClient *SFClient) SetDeviceId(deviceId string) *SFClient {
	sfClient.Api.Req.Settings.SetDeviceToken(deviceId)
	return sfClient
}
func (sfClient *SFClient) SetVersion(version string) *SFClient {
	sfClient.Api.Req.Settings.SetVersion(version)
	return sfClient
}
