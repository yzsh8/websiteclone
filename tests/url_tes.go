package test

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"testing"
)

/**
* 测试url
 */
func TestZhUrl(t *testing.T) {

	sHref := "https://shanghai.baixing.com/shangpuzhuanrang/style.css?page=1&cc=22#dd=33"

	oUrl, _ := url.ParseRequestURI(sHref)
	fmt.Println("ob->", oUrl.Path)
	fmt.Println("p1->", oUrl.RawQuery)
}

/**
* 测试正则
 */
func TestZhenzhe(t *testing.T) {
	str := `.i-case{width:100%;height:849px;background:url(../images/case_bj.png)no-repeatcentertop;position:relative;padding:73px0px0px0px;margin:50px0px;}
.i-caseulli{width:259px;background:#fff;border:1pxsolid#f0eded;padding:6px6px20px6px;float:left;text-align:center;margin-right:19px;margin-bottom:30px;}
.i-caseullispan{display:block;line-height:24px;height:24px;overflow:hidden;font-size:14px;margin-top:20px;}
.i-caseulliem{display:block;line-height:24px;color:#777;font-size:14px;}
.i-caseulli.a_lj{display:block;background:url(../images/jia.png)no-repeatcenter;height:30px;margin-top:10px;}
.i-caseulli.a_lj:hover{background:url(../images/jia_b.png)no-repeatcenter;}
.i-caseulli.no{margin-right:0px;}


.wrap{width:1150px;}
.a10{text-decoration:none;}
.md-head{line-height:58px;height:58px;background:url(../images/case_bt-bj.png)repeat-x;}
.md-heada{line-height:58px;float:left;height:58px;width:164px;font-size:16px;background:url(../images/fgx-s.png)repeat-yright;text-align:center;color:#fff;}
.md-heada.wid{width:164px}
.md-heada.cur{height:58px;line-height:58px;color:#ffffff;background:#ab8a56;}
.md-body{padding:30px0px;zoom:1;text-align:left;position:relative;}
.more{position:absolute;height:58px;width:150px;top:-58px;right:10px;}
.morea{background:url(../images/case-gd.png)no-repeatrightcenter;color:#bcbcbc;height:58px;line-height:58px;text-align:right;padding-right:35px;display:block;font-size:14px;}
.morea:hover{background:url(../images/case-gd_b.png)no-repeatrightcenter;color:#ab8a56;}



.our{width:1150px;}
.our_titel{background:url(../images/hxian.png)no-repeat100pxbottom;}
.our_titelstrong{display:block;font-size:24px;line-height:30px;color:#a8a8a8;font-weight:normal;}
.our_titelspan{display:block;float:left;line-height:18px;font-size:18px;color:#ab8a56;padding-top:16px;}
.our_titelp{text-align:right;color:#ab8a56;line-height:24px;height:24px;font-size:14px;padding-bottom:10px;}
.our_titela{margin-left:30px;}

.ourul{width:1150px;margin:0auto;margin-top:30px;}
.ourulli{float:left;width:215px;height:135px;margin-right:16px;margin-bottom:16px;text-align:center;border:1pxsolid#e5e5e5;}
.ourullispan{display:block;text-align:center;line-height:36px;font-size:18px;margin-top:10px;}
.ourullisamp{display:block;text-align:center;line-height:18px;font-size:14px;font-family:arial,helvetica,sans-serif;color:#999}
.ourulli.no{margin-right:0px;background:none;}


.fwlc{width:100%;height:666px;background:url(../images/fw_bj.png)no-repeatcentertop;margin:50px0px;}
.fwlc_bt{width:1150px;margin:0auto;padding-top:50px;padding-bottom:40px;}
.fwlc_btstrong{display:block;text-align:center;font-weight:normal;font-size:40px;color:#c9c8c8;}
.fwlc_btspan{margin-top:15px;display:block;line-height:30px;text-align:center;color:#c9c8c8;font-size:22px;background:url(../images/fw_x.png)no-repeatcentercenter;}
.lct{width:1180px;margin:0auto;text-align:center;color:#999;line-height:25px;}
.lctulli{float:left;width:235px;height:208px;text-align:center;border-right:1pxsolid#717273;}
.lctullia{display:block;text-align:center;color:#dfdfdf;transition:background1s;padding:34px0px;}
.lctulliimg{border-radius:41px;border:2pxsolid#ccc;}
.lctullia:hoverimg{background:#0171b8;border:2pxsolid#0171b8;}
.lctulliaspan{display:block;text-align:center;line-height:36px;font-size:14px;margin-top:10px;}
.lctulli.border{border-bottom:1pxsolid#717273;}
`
	flysnowRegexp := regexp.MustCompile(`url\(\"\S*?\"`)
	cmps := flysnowRegexp.FindStringSubmatch(str)

	for _, cmp := range cmps {
		cmp = strings.ReplaceAll(cmp, `url(`, "")
		cmp = strings.ReplaceAll(cmp, `"`, "")
		fmt.Println(cmp)
	}

	matched, err := regexp.MatchString(`url\(\"\S*?\"\)`, str)
	fmt.Println(matched, err)

	for i := 0; i <= 3; i++ {
		fmt.Println("run")
	}
}
