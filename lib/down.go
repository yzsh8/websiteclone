package lib

import (
	"errors"
	"time"
	"websiteclone/helper"
)

/**
* 下载主结构体
 */
type DOWN struct {
	PL  int64  //下载目录深度
	TO  int64  //超时重试时间
	PR  string //文件保存路径
	URL string //准备下载的网站地址
	CT  string //字符串编码
}

/**
* 获取html内容并保存到index.html中
 */
func (D *DOWN) Do() (bool, error) {
	res, l, n := HttpRequestByHeaderFor5(D.URL, "GET", "", nil)
	if n != 200 {
		time.Sleep(time.Duration(D.TO) * time.Second)
		res, l, n = HttpRequestByHeaderFor5(D.URL, "GET", "", nil)
	}

	if n != 200 {
		//下载失败
		return false, errors.New("下载文件失败:" + D.URL)
	}

	//如果状态没有返回，则计算字符串长度
	if l == 0 {
		l = int64(len(res))
	}

	//如果本地文件已经存在，则通过判断大小累决定是否需要更新
	sLocal := D.PR + "/index.html"
	cu := D.ChkUpdate(sLocal, l)

	//不需要更新
	if cu == false {
		//文件不需要更新
		return true, nil
	}

	//更新写入本地文件
	helper.CreateFile(sLocal) //创建目录
	_, er := helper.WriteFile(D.PR, "index.html", res)
	if er != nil {
		//文件写入失败
		return false, errors.New("写入本地文件失败:" + sLocal)
	}
	return true, nil
}

/**
* 判断本地文件是否需要更新
 */
func (D *DOWN) ChkUpdate(path string, slen int64) bool {
	if helper.IsExist(path) == true {
		localLen, _ := helper.GetFileSize(path)
		if slen == localLen {
			return false
		} else {
			return true
		}
	}
	return true
}
