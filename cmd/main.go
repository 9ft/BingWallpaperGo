package main

import (
	"BingWallpaperGo"
	"fmt"
)

func main() {
	host := "cn.bing.com"
	path := "path"
	err := BingWallpaperGo.SetBingWp(host, path)
	fmt.Println(err)
}
