package lib

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/simplejia/namecli/api"
	"github.com/simplejia/utils"
)

func TestPost(h http.HandlerFunc, params interface{}) (body []byte, err error) {
	v, err := json.Marshal(params)
	if err != nil {
		return
	}
	r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(v))
	if err != nil {
		return
	}
	w := httptest.NewRecorder()
	h(w, r)
	body = w.Body.Bytes()
	if g, e := w.Code, http.StatusOK; g != e {
		err = fmt.Errorf("http resp status not ok: %s", http.StatusText(g))
		return
	}
	return
}

func DeduplicateInt64s(a []int64) (result []int64) {
	if len(a) == 0 {
		return
	}

	exists := map[int64]bool{}

	result = a[:0]
	for _, e := range a {
		if exists[e] {
			continue
		}
		exists[e] = true

		result = append(result, e)
	}

	return
}

func SearchInt64s(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func BingoDisorderInt64s(a []int64, x int64) bool {
	for _, e := range a {
		if e == x {
			return true
		}
	}

	return false
}

func Int64s(a []int64) {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
}

func ZipInt64s(a []int64) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	err = json.NewEncoder(zw).Encode(a)
	if err != nil {
		return
	}
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipInt64s(a []byte) (result []int64, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	err = json.NewDecoder(zr).Decode(&result)
	if err != nil {
		return
	}
	return
}

func ZipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	zw.Write(a)
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	bs, err := ioutil.ReadAll(zr)
	if err != nil {
		return
	}

	result = bs
	return
}

func NameWrap(name string) (addr string, err error) {
	if strings.HasSuffix(name, ".ns") {
		return api.Name(name)
	}

	return name, nil
}

func PostProxy(name, path string, req []byte) (rsp []byte, err error) {
	addr, err := NameWrap(name)
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/%s", addr, strings.TrimPrefix(path, "/"))

	gpp := &utils.GPP{
		Uri:    url,
		Params: req,
	}
	rsp, err = utils.Post(gpp)
	if err != nil {
		return
	}

	return
}

func PostProxyReturnHeader(name, path string, req []byte) (rsp []byte, header http.Header, err error) {
	addr, err := NameWrap(name)
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/%s", addr, strings.TrimPrefix(path, "/"))

	reader := bytes.NewReader(req)
	r, err := http.Post(url, "application/json", reader)
	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		return
	}
	rsp, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if g, e := r.StatusCode, http.StatusOK; g != e {
		err = fmt.Errorf("http resp code: %d", g)
		return
	}

	header = r.Header
	return
}

func PostProxyWithHeader(name, path string, req []byte, headers map[string]string) (rsp []byte, err error) {
	addr, err := NameWrap(name)
	if err != nil {
		return
	}
	url := fmt.Sprintf("http://%s/%s", addr, strings.TrimPrefix(path, "/"))

	gpp := &utils.GPP{
		Uri:     url,
		Params:  req,
		Headers: headers,
	}
	rsp, err = utils.Post(gpp)
	if err != nil {
		return
	}

	return
}

func ClientWithProxy(name string) (client *http.Client, err error) {
	if name == "" {
		client = &http.Client{}
		return
	}

	addr, err := NameWrap(name)
	if err != nil {
		return
	}

	client = &http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse(fmt.Sprintf("http://%s", addr))
			},
		},
	}
	return
}

func TrimDataURL(data string) (ret string) {
	if data == "" {
		return
	}

	return regexp.MustCompile(`^data:[0-9a-zA-Z/]+?;base64,`).ReplaceAllString(data, "")
}

func TruncateWithSuffix(data string, length int, suffix string) (ret string) {
	rdata := []rune(data)
	if len(rdata) > length {
		ret = string(rdata[:length]) + suffix
	} else {
		ret = data
	}

	return
}
