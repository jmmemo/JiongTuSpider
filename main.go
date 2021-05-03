package main

import (
	"fmt"
	"github.com/jackdanger/collectlinks"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
var (
	URL   = "https://www.gamersky.com/ent/202105/1384885.shtml" //要爬取的链接
	Page  = 42                                                  //要爬取多少页
	Count int
)
*/

var (
	Count int
)

type Config struct {
	Name string
}

func (c *Config) InitConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	return viper.ReadInConfig()
}

func GetSomeInterestingPic(url string, page int) {
	var firstFlag = true
	//var linkArr []string   //wrong place

	split := strings.Split(url, ".")
	//fmt.Println(split)
	lens := len(split)
	s := split[lens-2] //slice to replace
	var homeRes []string
	for i := 1; i < page+1; i++ { // decide how many pages
		temp := split
		tmp := fmt.Sprintf("%s_%s", s, strconv.Itoa(i)) //Splicing pageNum
		//fmt.Println("数字拼接成功", tmp)
		//拼接切片
		temp[len(temp)-2] = tmp
		homeRes = append(homeRes, strings.Join(temp, "."))
	}
	//fmt.Printf("%#v\n", homeRes)

	for _, re := range homeRes {
		var linkArr []string
		if firstFlag {
			re = url
			//fmt.Println("第一页是", re)
			firstFlag = false
		}
		//对最外层链接遍历
		resp, err := http.Get(re)
		if err != nil {
			fmt.Println("http get error", err)
			return
		}

		links := collectlinks.All(resp.Body)
		for _, link := range links { //对最外层链接遍历，从里面寻找 pic link
			//fmt.Println(link)
			if strings.Index(link, "http") != 0 { //过滤非http链接
				//fmt.Println("Skip an abnormal link==>", link)
				continue
			}
			if strings.Index(link, "img1") != -1 {
				//fmt.Println("Real pic link-->", link)
				//save to a new []string => linkArr
				linkArr = append(linkArr, link)
			}
		} //now start range the linkArr
		for _, v := range linkArr {
			//split by "?" get the [1]
			RealPicUrl := strings.Split(v, "?")
			//fmt.Println("this v", RealPicUrl)
			resp, err := http.Get(RealPicUrl[1])
			if err != nil {
				panic(err)
			}
			data, err := ioutil.ReadAll(resp.Body) //get pic binary
			if err != nil && err != io.EOF {
				fmt.Println("read pic data failed,", err)
				return
			}
			str := fmt.Sprintf("./pic_%d.png", Count)
			file, err := os.OpenFile(str, os.O_CREATE|os.O_RDWR, 0644) //create file to write
			if err != nil {
				panic(err)
			}
			_, err = file.Write(data)
			if err != nil {
				fmt.Println("write err,", err)
				return
			}
			Count++
		}
	} //

}

func main() {
	c := Config{}
	err := c.InitConfig()
	if err != nil {
		panic(err)
	}
	URL := viper.Get("URL").(string)
	Page := viper.Get("Page").(int)

	startTime := time.Now()
	//rand.Seed(time.Now().Unix())
	log.Printf("爬取链接:%s 爬取页数:%d", URL, Page)
	GetSomeInterestingPic(URL, Page)
	log.Printf("爬取 %d张 囧图", Count)
	fmt.Printf("Done...Spent %f seconds!", time.Since(startTime).Seconds())
}
