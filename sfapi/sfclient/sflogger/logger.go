package sflogger

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type SFLogger struct {
	req           *resty.Request
	rep           *resty.Response
	formatLogText string
}

func NewSFLogger(rep *resty.Response) *SFLogger {
	return &SFLogger{req: rep.Request, rep: rep}
}

func (sfRequestLogger *SFLogger) CreateLogInfo() string {
	return sfRequestLogger.formatRequestLogText() + sfRequestLogger.formatResponseLogText()
}
func (sfRequestLogger *SFLogger) formatRequestLogText() string {
	var reqLogText string
	rh := copyHeaders(sfRequestLogger.req.Header)

	var body string
	if sfRequestLogger.req.QueryParam != nil {
		body = sfRequestLogger.req.QueryParam.Encode()
	}
	reqLogText = formatLog("\n==============================================================================\n"+
		"~~~ REQUEST ~~~\n"+
		"HOST   : %v\n"+
		"HEADERS:\n%s\n"+
		"BODY   :\n%v\n"+
		"------------------------------------------------------------------------------\n",
		sfRequestLogger.rep.RawResponse.Request.Host,
		composeHeaders(rh), body,
	)

	return reqLogText
}

func (sfRequestLogger *SFLogger) formatResponseLogText() string {
	var repLogText string
	if cookies := sfRequestLogger.rep.RawResponse.Cookies(); cookies != nil {
		repLogText += "  Cookies:\n"
		for _, cookie := range sfRequestLogger.rep.Cookies() {
			repLogText += fmt.Sprintf("    %s=%s", cookie.Name, cookie.Value)
		}
	}
	repLogText += formatLog("\n\n"+
		"~~~ RESPONSE ~~~\n"+
		"%s: %s  %s\n"+
		"Code   : %v\n"+
		"Status : %s\n"+
		"HEADERS:\n%s\n"+
		"BODY   :\n%v\n"+
		"------------------------------------------------------------------------------\n",
		sfRequestLogger.rep.Request.Method, sfRequestLogger.rep.Request.RawRequest.URL.RequestURI(), sfRequestLogger.rep.Proto(),
		sfRequestLogger.rep.StatusCode(),
		sfRequestLogger.rep.Status(),
		composeHeaders(sfRequestLogger.rep.Header()),
		indentJson(sfRequestLogger.rep.String()))
	return repLogText
}

func (sfRequestLogger *SFLogger) handleCookies(rh http.Header) string {
	var cookieText string
	if sfRequestLogger.req.Cookies != nil {
		for _, cookie := range sfRequestLogger.req.Cookies {
			cookieText += fmt.Sprintf("%s=%s\n", cookie.Name, cookie.Value)
		}
	} else {
		if rh.Get("Cookie") != "" {
			for _, cookie := range rh["Cookie"] {
				cookieText += fmt.Sprintf("%s\n", cookie)
			}
		}
	}
	return strings.TrimSpace(cookieText)
}

func formatLog(format string, params ...interface{}) string {
	return fmt.Sprintf(format, params...)
}
