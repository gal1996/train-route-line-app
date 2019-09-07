package handler

import (
	"log"
	"net/url"
	"strconv"
	"time"
	"unicode/utf8"
)

func GetUrl(start string, target string) (string, error) {
	routeUrl, err := makeUrl(start, target)
	if err != nil {
		return "", nil
	}
	return routeUrl, nil
}

func makeUrl(start string, target string) (string, error) {
	baseUrl, err := url.Parse("https://transit.yahoo.co.jp/search/result")
	if err != nil {
		return "", err
	}
	t := time.Now()
	nowUtc := t.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJst := nowUtc.In(jst)
	// リクエストbodyの作成
	values := url.Values{}
	values.Add("flatlon", "")
	values.Add("fromgid", "")
	values.Add("from", start)
	values.Add("togid", "")
	values.Add("viacode", "")
	values.Add("via", "")
	values.Add("to", target)
	values.Add("y", strconv.Itoa(t.Year()))
	month := strconv.Itoa(int(t.Month()))
	if utf8.RuneCountInString(month) == 1 {
		month = "0"+month
	}
	values.Add("m", month)
	log.Printf("m:%s", values.Get("m"))
	values.Add("d", strconv.Itoa(nowJst.Day()))
	log.Printf("d:%s", values.Get("d"))
	values.Add("hh", strconv.Itoa(nowJst.Hour()))
	log.Printf("hh:%s", values.Get("hh"))
	// 分数は別々に扱うので、分割する
	min := strconv.Itoa(nowJst.Minute())
	var m1 string
	var m2 string
	if utf8.RuneCountInString(min) == 2 {
		m2 = min[1:]
		m1 = min[:1]
	} else {
		m2 = min
		m1 = "0"
	}
	values.Add("m1", m1)
	log.Printf("m1:%s", values.Get("m1"))
	values.Add("m2", m2)
	log.Printf("m2:%s", values.Get("m2"))
	values.Add("type", "1")
	values.Add("ticket", "ic")
	values.Add("expkind", "1")
	values.Add("ws", "3")
	values.Add("s", "0")
	values.Add("al", "1")
	values.Add("shin", "1")
	values.Add("ex", "1")
	values.Add("hb", "1")
	values.Add("lb", "1")
	values.Add("sr", "1")
	values.Add("kw", target)

	baseUrl.RawQuery = values.Encode()

	routeUrl := baseUrl.String()
	log.Printf("made url : %s", routeUrl)

	return routeUrl, nil
}
