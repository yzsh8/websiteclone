package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/**
* 将复杂请求，因为网络的不稳定，将常规请求变成5次连续请求
 */
func HttpRequestByHeaderFor5(sUrl, method, sParame string, mHeader map[string]string) (string, int64, int) {
	res, l, n := HttpContentByHeader(sUrl, sParame, mHeader, method)
	if n != 200 {
		for i := 0; i < 5; i++ {
			res, l, n = HttpContentByHeader(sUrl, sParame, mHeader, method)
			if n == 200 {
				break
			}
		}
	}
	return res, l, n
}

/**
* 将复杂请求带cookie，因为网络的不稳定，将常规请求变成5次连续请求
 */
func HttpRequestByCookieFor5(sUrl, method, sParame string, mHeader map[string]string, cook []*http.Cookie) (string, int, []*http.Cookie) {
	res, n, cookie := httpRequestByCookie(sUrl, method, sParame, mHeader, cook)
	if n != 200 {
		for i := 0; i < 5; i++ {
			res, n, cookie = httpRequestByCookie(sUrl, method, sParame, mHeader, cook)
			if n == 200 {
				break
			}
		}
	}
	return res, n, cookie
}

func httpGetImgFor5(sUrl, sPath string, setCook []*http.Cookie) (string, []*http.Cookie) {
	res, cook := httpGetImg(sUrl, sPath, setCook)
	if len(res) < 1 {
		for i := 0; i < 5; i++ {
			res, cook = httpGetImg(sUrl, sPath, setCook)
			if len(res) > 1 {
				break
			}
		}
	}
	return res, cook
}

/**
*  简单版http请求，适用于没有特别要求的
* httpUrl	请求的网址
* method	网络请求方式，一般为POST或者GET
* sParam	需要传递的参数
* mHeader	http的头部
 */
func httpRequest(httpUrl, method, sParam string, mHeader map[string]string) string {
	client := &http.Client{}

	req, er := http.NewRequest(method, httpUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {
		fmt.Println("http request error->", er.Error())

		req, er = http.NewRequest(method, httpUrl, strings.NewReader(sParam))
		if er != nil {
			//两次连接都失败了，需要返回一个空
			return ""
		}
	}
	defer req.Body.Close()

	for key, val := range mHeader {
		req.Header.Add(key, val)
	}

	var body []byte
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			httpRequest(req.RequestURI, req.Method, "", mHeader)
		}

		body, _ = ioutil.ReadAll(resp.Body)
	} else {

		fmt.Println("http clident do error->", err.Error())
	}

	// fmt.Println("url->", httpUrl)
	// fmt.Println("参数->", sParam)
	// fmt.Println("返回->", string(body))

	return string(body)
}

/**
* 复杂版http，请求携带cookie
* httpUrl	请求的网址
* method	网络请求方式，一般为POST或者GET
* sParam	需要传递的参数
* mHeader	http的头部
* setCookie	传递的cookie
 */
func httpRequestByCookie(httpUrl, method, sParam string, mHeader map[string]string, setCookie []*http.Cookie) (string, int, []*http.Cookie) {
	src := ""
	httpStart := true
	statusCode := 101

	cook := []*http.Cookie{}

	req, er := http.NewRequest(method, httpUrl, bytes.NewReader([]byte(sParam)))
	if er != nil {
		fmt.Println("http request error->", er.Error())

		req, er = http.NewRequest(method, httpUrl, strings.NewReader(sParam))
		if er != nil {
			httpStart = false
			//两次连接都失败了，需要返回一个空
			return "", statusCode, setCookie
		}
	}

	for key, val := range mHeader {
		req.Header.Add(key, val)
	}

	if setCookie != nil && len(setCookie) > 0 {
		for _, v := range setCookie {
			req.AddCookie(v)
		}
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		if len(mHeader) > 0 {
			//mid := ""
			for key, val := range mHeader {
				//mid = key + ":" + val
				req.Header.Set(key, val)
			}
		}

		req.Header.Set("Accept-Charset", "utf-8")
		//req.Header.Set("Connection", "Close")

		tr := &http.Transport{DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			}}
		client := &http.Client{Transport: tr}
		req2 := req
		resp, err := client.Do(req2)

		if err != nil {
			fmt.Println("连接失败", err.Error())
			//两次连接都失败了，需要返回一个空
			return "", statusCode, setCookie
			fmt.Println("err->", httpUrl+err.Error())
		} else {
			defer resp.Body.Close()

			statusCode = resp.StatusCode
			cook = resp.Cookies()
			contents, _ := ioutil.ReadAll(resp.Body)
			src = string(contents)
		}
	}

	defer req.Body.Close()

	// fmt.Println("url->", httpUrl)
	// fmt.Println("参数->", sParam)
	// fmt.Println("返回->", string(body))

	return src, statusCode, cook
}

/**
* 复杂版http请求
* @param 	string 	sUrl	请求的地址
* @param	string	params	带入的参数
* @param	map	mHeader		head头
* @param	string	method		http方法
* 返回 字符串，长度，状态码
 */
func HttpContentByHeader(sUrl, params string, mHeader map[string]string, method string) (string, int64, int) {
	var strLen int64 = 0
	src := ""
	httpStart := true
	statusCode := 101
	req, err := http.NewRequest(method, sUrl, strings.NewReader(params))
	if err != nil {

		fmt.Println(sUrl + err.Error())

		req, err = http.NewRequest(method, sUrl, strings.NewReader(params))
		if err != nil {
			httpStart = false
			//两次连接都失败了，需要返回一个空
			return "", strLen, statusCode
		}
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		if len(mHeader) > 0 {
			//mid := ""
			for key, val := range mHeader {
				//mid = key + ":" + val
				req.Header.Set(key, val)
			}
		}

		req.Header.Set("Accept-Charset", "utf-8")
		req.Header.Set("Connection", "Close")

		tr := &http.Transport{DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			}}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(sUrl + err.Error())
		} else {

			statusCode = resp.StatusCode
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(sUrl + err.Error())
			} else {
				src = string(contents)
			}

		}
		if resp != nil && resp.Body != nil {
			strLen = resp.ContentLength
			defer resp.Body.Close()
		}
	}

	// fmt.Println("url->", sUrl)
	// fmt.Println("参数->", params)
	// fmt.Println("返回->", src)

	return src, strLen, statusCode
}

/**
* 此函数包含baixing网专属的header
* @param 	string 	sUrl	请求的地址
* @param	string	params	带入的参数
* @param	map	mHeader		head头
* @param	string	method		http方法
* setCookie	传递的cookie
 */
func HttpContentByHeaderByBaixing(sUrl, params string, mHeader map[string]string, method string, city string) (string, int) {
	src := ""
	httpStart := true
	statusCode := 101
	req, err := http.NewRequest(method, sUrl, strings.NewReader(params))
	if err != nil {

		fmt.Println(sUrl + err.Error())

		req, err = http.NewRequest(method, sUrl, strings.NewReader(params))
		if err != nil {
			httpStart = false
			//两次连接都失败了，需要返回一个空
			return "", statusCode
		}
	}

	if req != nil && req.Body != nil {
		defer req.Body.Close()
	}

	//只有连接成功后，才会写入头的读取字节流
	if httpStart == true {
		//模拟一个host
		sHost := ""
		u, err := url.Parse(sUrl)

		if err != nil {
			fmt.Println(sUrl + err.Error())
		} else {
			sHost = u.Host
		}

		if len(mHeader) > 0 {
			//mid := ""
			for key, val := range mHeader {
				//mid = key + ":" + val
				req.Header.Set(key, val)
			}
		}

		req.Header.Set("Cookie", "Hm_lpvt_5a727f1b4acc5725516637e03b07d3d2=1609942755; Hm_lvt_5a727f1b4acc5725516637e03b07d3d2=1609865217,1609934863; __sense_session_pv=22; _ga=GA1.2.1689221819.1609865216; _gid=GA1.2.1337414640.1609865216; _gat=1; _auth_redirect=https%3A%2F%2F"+city+".baixing.com%2Fshangpuzhuanrang%2F; __s=339mnmi2c8ga71kqa9utruftl7; __city="+city+"; __trackId=160986521481166")
		req.Header.Set("Origin", u.Scheme+"://"+u.Host)
		req.Header.Set("Host", sHost)
		req.Header.Set("Referer", sHost)
		req.Header.Set("Accept-Charset", "utf-8")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.1 Safari/605.1.15")
		req.Header.Set("Content-Type", "application/xml; charset=UTF-8")
		req.Header.Set("Connection", "Close")

		tr := &http.Transport{DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*30) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(5 * time.Second)) //设置发送接收数据超时
				return c, nil
			}}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(sUrl + err.Error())
		} else {

			statusCode = resp.StatusCode
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(sUrl + err.Error())
			} else {
				src = string(contents)
			}

		}
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
	}
	return src, statusCode
}

/**
*下载验证码
* httpUrl	请求的网址
* fileName	保存到本地的文件名
* setCookie	传递的cookie
 */
func httpGetImg(httpUrl, fileName string, setCookie []*http.Cookie) (string, []*http.Cookie) {
	client := &http.Client{}
	cook := []*http.Cookie{}

	req, er := http.NewRequest("GET", httpUrl, nil)
	if er != nil {
		fmt.Println("get img err->", er.Error())
		return "", cook
	}
	req.Close = true
	if setCookie != nil && len(setCookie) > 0 {
		for _, v := range setCookie {
			req.AddCookie(v)
		}
	}

	fmt.Println("img url->", httpUrl)

	var body []byte
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read captcha->", err.Error())
		}
		//读取最新的cook
		cook = resp.Cookies()

		out, er := os.Create(fileName)
		io.Copy(out, bytes.NewReader(body))

		if er != nil {
			fmt.Println("save captcha->", er.Error())
		}

		return fileName, cook
	}

	fmt.Println("get captcha error->", err.Error())

	return "", cook
}
