package lib

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"websiteclone/helper"

	"github.com/PuerkitoBio/goquery"
)

/**
* 下载主结构体
 */
type DOWN struct {
	PL  int64  //下载目录深度
	TO  int64  //超时重试时间
	PR  string //文件保存路径
	URL string //准备下载的网站地址
	DM  string //准备下载的网站的域名
}

/**
* 获取html内容并保存到index.html中
 */
func (D *DOWN) Do() (bool, error) {
	res, _, n := HttpRequestByHeaderFor5(D.URL, "GET", "", nil)
	if n != 200 {
		time.Sleep(time.Duration(D.TO) * time.Second)
		res, _, n = HttpRequestByHeaderFor5(D.URL, "GET", "", nil)
	}

	if n != 200 {
		//下载失败
		return false, errors.New("下载文件失败:" + D.URL)
	}

	//本地文件
	sLocal := D.PR + "/index.html"

	//读取html中的node节点，下载资源文件
	res = D.downTasks(res)

	//更新写入本地文件
	helper.CreateFile(D.PR) //创建目录
	_, er := helper.WriteFile(D.PR, "index.html", res)
	if er != nil {
		//文件写入失败
		return false, errors.New("写入本地文件失败:" + sLocal)
	}

	return true, nil
}

/**
* 下载资源文件的任务
* 返回更新的html（需要修改远程css和js到本地来）
 */
func (D *DOWN) downTasks(str string) string {
	//字符串全部转成小写
	str = strings.ToLower(str)

	nHead, _ := getHead(str)
	heads := goquery.NewDocumentFromNode(nHead)

	//获取head头中需要下载的link列表
	heads.Find("link").Each(func(i int, s *goquery.Selection) {
		sHref, _ := s.Attr("href")

		//下载
		sLocalUrl := D.downPool(sHref)
		str = strings.Replace(str, sHref, sLocalUrl, 0)
	})

	//获取head头中需要下载的script文件
	heads.Find("script").Each(func(i int, s *goquery.Selection) {
		sHref, _ := s.Attr("src")

		//下载
		sLocalUrl := D.downPool(sHref)
		str = strings.Replace(str, sHref, sLocalUrl, 0)
	})

	//获取body中需要下载的文件
	nBody, _ := getBody(str)
	bodys := goquery.NewDocumentFromNode(nBody)

	//下载js
	bodys.Find("script").Each(func(i int, s *goquery.Selection) {
		sHref, _ := s.Attr("src")
		sLocalUrl := D.downPool(sHref)
		str = strings.Replace(str, sHref, sLocalUrl, 0)
	})

	//下载css
	bodys.Find("link").Each(func(i int, s *goquery.Selection) {
		sHref, _ := s.Attr("href")
		sLocalUrl := D.downPool(sHref)
		str = strings.Replace(str, sHref, sLocalUrl, 0)
	})

	//下载图片
	bodys.Find("img").Each(func(i int, s *goquery.Selection) {
		sHref, _ := s.Attr("src")
		sLocalUrl := D.downPool(sHref)
		str = strings.Replace(str, sHref, sLocalUrl, 0)
	})

	return str
}

/**
*下载资源文件的子进程
 */
func (D *DOWN) downPool(sHref string) string {
	localUrl := ""

	fmt.Println("down pool ->", sHref)

	//空的不需要下载
	if len(sHref) < 1 {
		return sHref
	}

	//补全完整url,当url地址是相对地址的时候
	if strings.Contains(sHref, "http:") == false {
		//如果是//开头的
		if strings.Index(sHref, "//") == 0 {
			sHref = "http:" + sHref
		} else {
			if strings.Index(sHref, "/") == 0 {
				sHref = D.URL + sHref
			} else {
				sHref = D.URL + "/" + sHref
			}
		}
	}

	//判断是否跨域，如果跨域的，不需要下载
	if strings.Contains(sHref, D.DM) == false {
		return sHref
	}

	//获取需要下载的文件名称
	oUrl, _ := url.ParseRequestURI(sHref)
	loalPath := oUrl.Path
	localUrl = "/" + loalPath

	//获取参数
	localParam := oUrl.RawQuery
	if len(localParam) > 1 {
		localUrl = "?" + localParam
	}

	//获取远程文件的大小
	res, l, n := HttpRequestByHeaderFor5(sHref, "GET", "", nil)
	if n != 200 {
		//远程文件不存在
		return localUrl
	}

	//如果是图片或者css文件，优先判断文件是否一样大，来决定是否重新下载
	if strings.Contains(sHref, ".css") == true || strings.Contains(sHref, ".png") == true || strings.Contains(sHref, ".gif") == true && strings.Contains(sHref, ".jpg") == true {
		localSize, _ := helper.GetFileSize(D.PR + loalPath)
		if localSize == l {
			//相同大小，不用下载
			fmt.Println("已经缓存->", sHref, l, localSize)
			return localUrl
		}
	}

	//下载到本地
	downFileFor5(sHref, D.PR+loalPath)

	//如果下载的文件是css，则需要下载里面包含的backgroud图片
	if strings.Contains(sHref, ".css") == true {
		//得到css文件路径
		cssPath := strings.Split(loalPath, "/")

		css := strings.ToLower(res)
		css = strings.ReplaceAll(css, " ", "") //去掉空格

		//通过正则查询,带双引号的情况下
		flysnowRegexp := regexp.MustCompile(`url\(\"\S*?\"`)
		cmps := flysnowRegexp.FindAllString(css, 0)
		D.downCssimg(cssPath, cmps)

		//通过正则查询，不带双引号的情况下
		flysnowRegexp2 := regexp.MustCompile(`url[\(]\S*?[\)]`)
		cmps2 := flysnowRegexp2.FindStringSubmatch(css)
		D.downCssimg(cssPath, cmps2)
	}

	return localUrl
}

/**
* 下载css文件中的图片
 */
func (D *DOWN) downCssimg(cssPath, cmps []string) {

	fmt.Println("正则结果-》", len(cmps), cmps)

	//需要下载的图片列表
	for _, cmp := range cmps {
		cmp = strings.ReplaceAll(cmp, `url(`, "")
		cmp = strings.ReplaceAll(cmp, `"`, "")
		cmp = strings.ReplaceAll(cmp, `)`, "")

		fmt.Println("css中得到图片地址->", cmp)

		//如果是携带完整url，跳过
		if strings.Contains(cmp, "http") == true {
			continue
		}

		//相对路径的，需要得到正确路径
		cssTrue := D.URL
		nP := strings.Count(cmp, "../")
		for i := 0; i < len(cssPath)-nP-1; i++ {
			cssTrue = cssTrue + "/" + cssPath[i]
		}
		cssTrue = cssTrue + "/" + strings.ReplaceAll(cmp, "../", "")

		fmt.Println("css中计算得到正确地址->", cssTrue)

		//下载这个图片
		tmUrl, _ := url.ParseRequestURI(cssTrue)
		downFileFor5(cssTrue, D.PR+tmUrl.Path)
	}
}
