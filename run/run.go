package runtimes

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"websiteclone/helper"

	"github.com/astaxie/beego"
)

var (
	pathlen, _    = beego.AppConfig.Int64("pathlen")    //mul深度
	downpool, _   = beego.AppConfig.Int64("downpool")   //网站并发下载线程数
	synctime, _   = beego.AppConfig.Int64("synctime")   //下一次同步检测时间间隔
	rumtimeout, _ = beego.AppConfig.Int64("rumtimeout") //超时重试时间间隔
	character     = beego.AppConfig.String("character") //字符串编码
)

/**
* 程序入口
 */
func runs() {
	//读取需要同步的网站列表
	fpri, _ := os.Open("conf/website.ini")
	defer fpri.Close()
	byPri, _ := ioutil.ReadAll(fpri)
	//解析json
	jOut := map[string]interface{}{}
	//解析json数据
	json.Unmarshal(byPri, &jOut)

	cityLists := map[string]string{}
	for k, v := range jOut { //遍历
		val := helper.UnknowToString(v) //得到最终结果
		cityLists[k] = val              //最终结果
	}
}
