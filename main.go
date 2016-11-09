package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"fmt"
	"encoding/json"
	"strings"
	"time"
)

//----------------------------------
// 重庆、天津、新疆时时彩开奖结果 － 开彩网
// 在线接口文档：http://face.apius.cn/?token=demo
//----------------------------------

// 查询结果
type Cxjg struct {
	Rows int `json:"rows"`
	Code string `json:"code"`
	Info string `json:"info"`
	Data []*Kjjg `json:"data"`
}

// 开奖结果
type Kjjg struct {
	Expect        string `json:"expect"`
	Opencode      string `json:"opencode"`
	Opentime      string `json:"opentime"`
	Opentimestamp int64 `json:"opentimestamp"`
}

type Jg struct {
	Qs string `json:"qs"`
	Sj string `json:"sj"`
	Ww string `json:"ww"`
	Qw string `json:"qw"`
	Bw string `json:"bw"`
	Sw string `json:"sw"`
	Gw string `json:"gw"`
	Q3 string `json:"q3"`
	Z3 string `json:"z3"`
	H3 string `json:"h3"`
}

var (
	token = "1DD4EF9DBA54131115D84281D5F9BD5F"
	cqssc []*Jg
	xjssc []*Jg
	tjssc []*Jg
	cqnum int64 = 0
	xjnum int64
	tjnum int64
)

func main() {
	// 开始执行任务
	go SyncTask();

	mux := http.NewServeMux()
	mux.HandleFunc("/cqssc", Cqssc)
	mux.HandleFunc("/xjssc", Xjssc)
	mux.HandleFunc("/tjssc", Tjssc)

	// http监听端口
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		Error(err.Error())
	}
}

// 获取重庆时时彩开奖数据
func Cqssc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, _ := json.Marshal(cqssc)
	cqnum = cqnum + 1
	Info(fmt.Sprintf("access cqssc data times: %d", cqnum))
	w.Write(data)
}

// 获取新疆时时彩开奖数据
func Xjssc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, _ := json.Marshal(xjssc)
	xjnum = xjnum + 1
	Info(fmt.Sprintf("access xjssc data times: %d", xjnum))
	w.Write(data)
}

// 获取天津时时彩开奖数据
func Tjssc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, _ := json.Marshal(tjssc)
	tjnum = tjnum + 1
	Info(fmt.Sprintf("access tjssc data times: %d", tjnum))
	w.Write(data)
}

// 自动获取彩票开奖结果
func SyncTask() {
	c := time.Tick(1 * 60 * time.Second)
	cqssc = SyncData("cqssc")
	xjssc = SyncData("xjssc")
	tjssc = SyncData("tjssc")
	for _ = range c {
		cqssc = SyncData("cqssc")
		xjssc = SyncData("xjssc")
		tjssc = SyncData("tjssc")
	}
}

// 彩票开奖结果查询
func SyncData(code string) (ssc []*Jg) {
	//请求地址
	apiURL := fmt.Sprintf("http://t.apiplus.cn/newly.do?token=%s&rows=50&format=json&code=%s", token, code)

	//发送请求
	data, err := Get(apiURL, nil)
	if err != nil {
		Error(fmt.Sprintf("get data from %s : %s", apiURL, err.Error()))
	} else {
		Info(fmt.Sprintf("get data from %s : %s", apiURL, string(data)))
		var cxjg Cxjg
		json.Unmarshal(data, &cxjg)
		for _, kjjg := range cxjg.Data {
			codes := strings.Split(kjjg.Opencode, ",")
			var jg Jg
			jg.Qs = kjjg.Expect
			jg.Sj = Substr(kjjg.Opentime, 11, 5)
			jg.Ww = codes[0]
			jg.Qw = codes[1]
			jg.Bw = codes[2]
			jg.Sw = codes[3]
			jg.Gw = codes[4]
			if jg.Ww == jg.Qw || jg.Ww == jg.Bw || jg.Qw == jg.Bw {
				jg.Q3 = "组三"
			} else {
				jg.Q3 = "组六"
			}
			if jg.Qw == jg.Bw || jg.Qw == jg.Sw || jg.Bw == jg.Sw {
				jg.Z3 = "组三"
			} else {
				jg.Z3 = "组六"
			}
			if jg.Bw == jg.Sw || jg.Bw == jg.Gw || jg.Sw == jg.Gw {
				jg.H3 = "组三"
			} else {
				jg.H3 = "组六"
			}
			ssc = append(ssc, &jg)
		}
	}
	return ssc
}

// get 网络请求
func Get(apiURL string, params url.Values) (rs[]byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		Error(err.Error())
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	if params != nil {
		Url.RawQuery = params.Encode()
	}
	resp, err := http.Get(Url.String())
	if err != nil {
		Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs[]byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 截取字符串，start 开始下标，length 截取长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func Info(str string) {
	fmt.Print(time.Now().Format("2006-01-02 15:04:05.000"), " INFO - ")
	fmt.Println(fmt.Sprintf(str))
}

func Error(str string) {
	fmt.Print(time.Now().Format("2006-01-02 15:04:05.000"), " ERROR - ")
	fmt.Println(fmt.Sprintf(str))
}
