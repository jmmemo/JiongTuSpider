package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	//url := "https://www.gamersky.com/showimage/id_gamersky.shtml?https://img1.gamersky.com/image2021/05/20210501_zy_red_164_97/gamersky_01origin_01_202151908D9.jpg"
	url := "https://img1.gamersky.com/image2021/05/20210501_zy_red_164_97/gamersky_01origin_01_202151908D9.jpg"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("pic1.png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
}
