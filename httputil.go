package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
	{
     "msgtype": "markdown",
     "markdown": {
         "title":"杭州天气",
         "text": "#### 杭州天气 @150XXXXXXXX \n> 9度，西北风1级，空气良89，相对温度73%\n> ![screenshot](https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png)\n> ###### 10点20分发布 [天气](https://www.dingtalk.com) \n"
     },
      "at": {
          "atMobiles": [
              "150XXXXXXXX"
          ],
          "isAtAll": false
      }
 }
*/

type Message struct {
	Msgtype  string  `json:"msgtype"`
	Markdown Content `json:"markdown"`
	ActionCard ActionCard `json:"actionCard"`
}
type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}



func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}


func sendMsg()  {
	dingUrl:=fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s",accessToken,timestamp,sign(timestamp,secret))
	bingBaseUrl := "https://cn.bing.com"

	bing:=getBingImg()
	imgInfo:=bing.Images[0]
	hdImage:=bingBaseUrl+imgInfo.Urlbase+HD
	uhdImage:=bingBaseUrl+imgInfo.Urlbase+UHD
	md:= "![](%s) \n> [%s](%s) " +
		"\n ---"+
		"\n # [下载壁纸（4k）](%s)"
	dingText:=fmt.Sprintf(md,hdImage,imgInfo.Copyright,imgInfo.Copyrightlink,uhdImage)
	fmt.Println(dingText)
	c:=Content{
		Title: "每日一图",
		Text:  dingText,
	}

	//bit1:=Btn{
	//	Title:     "下载壁纸(1920X1080)",
	//	ActionURL: hdImage,
	//}
	//bit2:=Btn{
	//	Title:     "下载壁纸(4k)",
	//	ActionURL: uhdImage,
	//}
	//
	//var actionCard = ActionCard{
	//	Title:          "每日一图",
	//	Text:           "![]("+hdImage+")",
	//	HideAvatar:     "0",
	//	BtnOrientation: "0",
	//	//SingleTitle: "下载",
	//	//SingleURL: "https://developers.dingtalk.com/document/app/custom-robot-access",
	//	Btns: 			[]Btn{bit2},
	//}
	msg:=Message{
		Msgtype: "markdown",
		//ActionCard: actionCard ,
		Markdown: c,

	}
	b,e:=json.Marshal(msg)
	if e!=nil {
		fmt.Println(e)
	}
	resp,e:=http.Post(dingUrl,"application/json",strings.NewReader(string(b)))
	defer resp.Body.Close()
	if e!=nil {
		fmt.Println(e)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString:=string(body)
	fmt.Println(bodyString)
}

func getBingImg() BingRep {
	bingUrl := "https://cn.bing.com/HPImageArchive.aspx/?format=js&n=1&uhd=1"
	resp,e:=http.Get(bingUrl)
	defer resp.Body.Close()
	if e!=nil {
		fmt.Println(e)
	}
	body,_:=ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var bingRep BingRep
	err:=json.Unmarshal(body, &bingRep)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println(bingRep.Images[0].Urlbase)
	return bingRep
}

const (

	HD = "_1920x1080.jpg"
	UHD = "_UHD.jpg"
)

type Img struct {
	Urlbase string `json:"urlbase"`
	Copyright string `json:"copyright"`
	Copyrightlink string `json:"copyrightlink"`
}
type BingRep struct {
	Images []Img `json:"images"`
}