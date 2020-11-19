package requests

import (
	"net/http"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

var session = grequests.NewSession(nil)

// * split each string of stringArray by sep
func parse(stringArray []string, sep string) map[string]string {
	res := map[string]string{}
	for _, s := range stringArray {
		temp := strings.SplitN(s, sep, 2)
		if len(temp) > 1 {
			res[temp[0]] = temp[1]
		} else {
			res[temp[0]] = ""
		}
	}
	return res
}

func Requests(_type string, url string, reqArgs map[string][]string, oArgs map[string]interface{}) (*grequests.Response, error) {
	isUseSession := false
	_type = strings.ToUpper(_type)
	resp := &grequests.Response{}
	var rerr error = errors.New("Unsupport type for reuqests")
	ro := &grequests.RequestOptions{}
	ro.RedirectLimit = 0
	// ? parse oArgs
	if r, ok := oArgs["UseSession"]; ok {
		isUseSession = r.(bool)
	}
	if r, ok := oArgs["UseSession"]; ok {
		isUseSession = r.(bool)
	}
	if r, ok := oArgs["Follow"]; ok && !r.(bool) {
		ro.HTTPClient = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	if r, ok := oArgs["ReqTimeout"]; ok {
		ro.RequestTimeout = time.Duration(r.(int)) * time.Second
	}
	if r, ok := oArgs["connDelayTimeout"]; ok {
		ro.TLSHandshakeTimeout = time.Duration(r.(int)) * time.Second
	}
	// ? parse reqArgs
	// parse headers
	if r, ok := reqArgs["Headers"]; ok {
		ro.Headers = parse(r, ":")
	}
	// parse cookies
	if r, ok := reqArgs["Cookies"]; ok && len(r) > 0 {
		cookies := []*http.Cookie{}
		for _, cookieString := range r {
			cookie := http.Cookie{}
			cookieArray := strings.SplitN(cookieString, "=", 2)
			cookie.Name = cookieArray[0]
			if len(cookieArray) > 1 {
				cookie.Value = cookieArray[1]
			} else {
				cookie.Value = ""
			}
			cookies = append(cookies, &cookie)
		}

		ro.Cookies = cookies
	}
	// parse data
	if rawData, ok := reqArgs["Data"]; ok && len(rawData) > 0 {
		// if data exists, force send post requests
		if _type == "GET" {
			_type = "POST"
		}
		ro.Data = parse(rawData, "=")

	}
	// parse json data
	if jsonData, ok := oArgs["Json"]; ok && len(jsonData.(string)) > 0 {
		if _type == "GET" {
			_type = "POST"
		}
		ro.Headers["Content-Type"] = "application/json"
		jsonData := jsonData.(string)
		ro.JSON = jsonData
	}
	// parse auth
	if r, ok := reqArgs["Auth"]; ok && len(r) > 0 {
		ro.Auth = r
	}
	// ? requests
	if isUseSession {
		switch strings.ToUpper(_type) {
		case "GET":
			resp, rerr = session.Get(url, ro)
		case "POST":
			resp, rerr = session.Post(url, ro)
		case "HEAD":
			resp, rerr = session.Head(url, ro)
		case "PUT":
			resp, rerr = session.Put(url, ro)
		case "DELETE":
			resp, rerr = session.Delete(url, ro)
		case "OPTIONS":
			resp, rerr = session.Options(url, ro)
		default:
		}
	} else {
		switch strings.ToUpper(_type) {
		case "GET":
			resp, rerr = grequests.Get(url, ro)
		case "POST":
			resp, rerr = grequests.Post(url, ro)
		case "HEAD":
			resp, rerr = grequests.Head(url, ro)
		case "PUT":
			resp, rerr = grequests.Put(url, ro)
		case "DELETE":
			resp, rerr = grequests.Delete(url, ro)
		case "OPTIONS":
			resp, rerr = grequests.Options(url, ro)
		default:

		}
	}
	// ! test
	// fmt.Println(resp.String())
	return resp, rerr
}
