package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//如果后面跟了参数，直接进行命令行操作
	if len(os.Args) > 1 {
		cliDo()
		fmt.Println("任务执行完成，正常退出")
		return
	}
}

/**
* 执行命令行操作
 */
func cliDo() {

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "websiteclone"
	app.Usage = "参数"
	app.UsageText = "本程序专门为克隆网站使用"
	app.ArgsUsage = ``

	app.Email = "yzsh8.james@gmail.com"
	app.Author = "cook"

	//这个方法就是这个命令已启动会运行什么
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c) //这个是打印app的help界面
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "runtimes",
			Aliases: []string{"run"},
			Usage: `启动计划任务:
					tuijiewang 		--推介网,[ta，tp],[验证码概率不高]
					`,
			Action: func(c *cli.Context) error {
				website := c.Args().First()

				switch website {
				case "tuijiewang":

				default:
					fmt.Println("输入错误，请使用-h查看帮助")
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
