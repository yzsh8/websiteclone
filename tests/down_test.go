package test

import (
	"fmt"
	"websiteclone/lib"

	"testing"
)

/**
* 测试url
 */
func TestDownUrl(t *testing.T) {
	do := new(lib.DOWN)
	do.URL = "http://www.yeeshone.com"
	do.DM = "www.yeeshone.com"
	do.PL = 1
	do.TO = 100
	do.PR = "./yeeshone"

	b, er := do.Do()

	if er != nil {
		fmt.Println("er->", er.Error())
	}

	fmt.Println("bug->", b)
}
