package sfclient

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func generateTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
}
func (s *SFRequest) sfSecurity() string {
	var (
		newMd5            = md5.New()
		sTimeStamp        = generateTimestamp()
		randomDeviceToken = s.Settings.GetRandomDeviceToken()
	)
	newMd5.Write([]byte(randomDeviceToken + sTimeStamp + s.Settings.GetDeviceToken() + s.Settings.GetAndroidKey()))
	sign := strings.ToUpper(hex.EncodeToString(newMd5.Sum(nil)))
	return fmt.Sprintf("nonce=%v&timestamp=%v&devicetoken=%v&sign=%v", randomDeviceToken, sTimeStamp, s.Settings.GetDeviceToken(), sign)
}

func (s *SFRequest) getHeaders() map[string]any {
	h := map[string]any{
		"Accept-Charset": "UTF-8",
		"Connection":     "Keep-Alive",
		"Accept":         "application/vnd.sfacg.api+json;version=1",
		"SFSecurity":     s.sfSecurity(),
		"User-Agent":     s.Settings.GetUserAgent(),
		"Content-Type":   "application/json",
	}
	return h
}
