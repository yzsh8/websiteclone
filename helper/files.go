package helper

import (
	"fmt"
	"io/ioutil"
	"os"
)

/**
* 创建文件,包含多层目录
 */
func CreateFile(src string) (string, error) {
	//	src := dir + name + "/"
	if IsExist(src) {
		return src, nil
	}

	if err := os.MkdirAll(src, 0777); err != nil {
		if os.IsPermission(err) {
			fmt.Println("你不够权限创建文件")
		}
		return "", err
	}

	return src, nil
}

/*
* 判断文件或者目录是否存在
 */
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/**
* 新建文件并写入内容
* 如果文件已存在,则覆盖以前内容
 */
func WriteFile(filePath, fileName, content string) (int, error) {
	_, err := CreateFile(filePath)
	if err != nil {
		return 0, err
	}
	src := filePath + "/" + fileName
	fs, e := os.Create(src)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

/**
* 获取文件大小,单位时B
 */
func GetFileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

/**
* 读取文件内容并返回字符串
* @param  path  文件路径
 */
func ReadFileString(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		//panic(err)
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
