/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-23 17:23:30
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-23 17:58:09
 * @FilePath: /goabao/abaozipfile/abaozipfile_test.go
 * @Description:
 */
package abaozipfile

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAbaoZipFile(t *testing.T) {
	//oldFileName可以是文件或者目录
	oldFileName := "test.log"

	currentTime := time.Now()

	//获取s
	mSecond := fmt.Sprintf("%03d", currentTime.Nanosecond()/1e6)

	//zip文件名
	zipFileName := strings.Split(oldFileName, ".")[0] + "_" + currentTime.Format("20060102150405") + mSecond + ".zip"

	//压缩文件
	err := CompressZip(oldFileName, zipFileName)
	if err != nil {
		t.Errorf("%+v", err)
	}

	err = DeCompress(zipFileName, "test")
	if err != nil {
		t.Errorf("%+v", err)
	}

}
