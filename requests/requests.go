package requests

import (
	"strings"

	"github.com/asmcos/requests"
	"github.com/pkg/errors"
)

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

func Requests(_type string, url string, args map[string][]string) (*requests.Response, error) {
	Nargs := make([]interface{}, 0)
	resp := &requests.Response{}
	var rerr error = errors.New("Unsupport type for reuqests")
	data := requests.Datas{}
	headers := requests.Header{}
	// ? parse
	// parse data
	if r, ok := args["Data"]; ok && len(r) > 0 {
		// if data exists, force send post requests
		_type = "POST"
		data = requests.Datas(parse(r, "="))
		Nargs = append(Nargs, data)

	}
	// parse headers
	if r, ok := args["Headers"]; ok {
		headers = requests.Header(parse(r, ":"))
	}
	// parse cookies
	if r, ok := args["Cookies"]; ok && len(r) > 0 {
		headers["Cookies"] = strings.Join(r, "&")
	}
	// parse cookies
	if r, ok := args["Auth"]; ok && len(r) > 0 {
		auth := requests.Auth(r)
		Nargs = append(Nargs, auth)
	}
	if len(headers) > 0 {
		Nargs = append(Nargs, headers)
	}

	// ? requests
	_type = strings.ToUpper(_type)
	if _type == "GET" {
		resp, rerr = requests.Get(url, Nargs...)
	}
	if _type == "POST" {
		resp, rerr = requests.Post(url, Nargs...)
	}
	// ! test
	// fmt.Println(resp.Text())
	return resp, rerr
}
