package sfsettings

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const defaultVersion = "4.8.42"
const defaultAndroidKey = "FMLxgOdsfxmN!Dt4"
const defaultWebBaseAPI = "https://pages.sfacg.com"
const defaultWebBaseURL = "https://book.sfacg.com"
const defaultAndroidBaseAPI = "https://api.sfacg.com"
const userAgentTemplate = "boluobao/%v(android;25)/HomePage/%v/HomePage"

const deviceTokenPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

type Settings struct {
	androidKey  string
	userAgent   string
	version     string
	cookies     string
	deviceToken string
	baseURL     string
}

func NewSettings() *Settings {
	settings := &Settings{
		version:    defaultVersion,
		androidKey: defaultAndroidKey,
		baseURL:    defaultAndroidBaseAPI,
	}
	return settings.SetDeviceToken(strings.ToLower(settings.GetRandomDeviceToken()))
}
func (s *Settings) SetVersion(version string) *Settings {
	s.version = version
	return s
}

func (s *Settings) SetUserAgent(userAgent string) *Settings {
	s.userAgent = userAgent
	return s
}
func (s *Settings) SetCookie(cookie any) *Settings {
	switch cookieInfo := cookie.(type) {
	case string:
		if strings.Contains(cookieInfo, "session_APP") && strings.Contains(cookieInfo, ".SFCommunity") {
			s.cookies = cookieInfo
		} else {
			log.Println("cookie type error")
		}
	case []*http.Cookie:
		var cookieStr string
		for _, i := range cookieInfo {
			cookieStr += i.Name + "=" + i.Value + ";"
		}
		s.SetCookie(cookieStr)
	default:
		log.Println("cookie type error")
	}
	return s
}
func (s *Settings) SetAndroidKey(androidKey string) *Settings {
	s.androidKey = androidKey
	return s
}

func (s *Settings) SetDeviceToken(deviceToken string) *Settings {
	if regexp.MustCompile(deviceTokenPattern).MatchString(deviceToken) {
		s.deviceToken = deviceToken
	} else {
		log.Println("deviceToken type error")
	}
	return s
}
func (s *Settings) SetBaseURL(baseURL string) *Settings {
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		log.Println("baseURL type error")
	} else {
		s.baseURL = baseURL
	}
	return s
}
func (s *Settings) GetVersion() string {
	return s.version
}
func (s *Settings) GetCookie() string {
	return s.cookies
}
func (s *Settings) GetDeviceToken() string {
	return strings.ToUpper(s.deviceToken)
}
func (s *Settings) GetRandomDeviceToken() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return strings.ToUpper(fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]))
}

func (s *Settings) GetAndroidKey() string {
	return s.androidKey
}
func (s *Settings) GetUserAgent() string {
	if s.userAgent == "" {
		return fmt.Sprintf(userAgentTemplate, s.version, s.deviceToken)
	}
	return s.userAgent
}
func (s *Settings) GetBaseAPI() string {
	return s.baseURL
}
func (s *Settings) GetWebBaseAPI() string {
	return defaultWebBaseAPI
}

func (s *Settings) GetWebBaseURL() string {
	return defaultWebBaseURL
}
