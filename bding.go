package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/greycodee/dingbot"
	"github.com/greycodee/dingbot/message"
	"io/ioutil"
	"net/http"
	"net/url"
)

var secret = flag.String("secret","","加签秘钥")
var accessToken = flag.String("token","","token")
func main()  {
	flag.Parse()
	sendMsg()
	
}

func sendMsg()  {
	fmt.Println(*secret)
	bot:=dingbot.DingBot{
		Secret:      *secret,
		AccessToken: *accessToken,
	}
	imgInfo := getBingImg().Images[0]
	img:=fmt.Sprintf("![](https://cn.bing.com%s)",imgInfo.Urlbase+HD)
	fmt.Println(img)
	imgUHD := fmt.Sprintf("https://cn.bing.com%s",imgInfo.Urlbase+UHD)
	imgUHD = url.QueryEscape(imgUHD)
	fmt.Println(imgUHD)
	msg:=message.Message{
		MsgType:   	message.ActionCardStr,
		ActionCard: message.ActionCard_{
			Title:          "必应每日一图",
			Text:           img+" \n ### 必应每日一图 \n "+imgInfo.Copyright,
			BtnOrientation: "0",
			HideAvatar:     "1",
			BtnS:           []message.Btn_{
				{
					Title:     "下载壁纸(4K)",
					ActionURL: fmt.Sprintf("dingtalk://dingtalkclient/page/link?url=%s&pc_slide=false",imgUHD),
				},
				{
					Title:     "copyright",
					ActionURL: imgInfo.Copyrightlink,
				},
			},
		},
	}
	bot.Send(msg)
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