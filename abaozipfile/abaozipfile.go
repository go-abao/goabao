/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-23 17:22:51
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-23 17:57:46
 * @FilePath: /goabao/abaozipfile/abaozipfile.go
 * @Description:
 */
package abaozipfile

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/**
 * @description: 压缩文件
 * @param {*} source
 * @param {string} target
 * @return {*}
 */
func CompressZip(source, target string) error {

	//创建目标zip文件
	zipFile, err := os.Create(target)

	if err != nil {
		fmt.Println(err)
		return err
	}

	//关闭文件
	defer zipFile.Close()

	//创建一个写zip的writer
	archive := zip.NewWriter(zipFile)

	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		//将文件或者目录信息转换为zip格式的文件信息
		header, err := zip.FileInfoHeader(info)

		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 确定采用的压缩算法（这个是内建注册的deflate）
			header.Method = zip.Deflate
		}

		//
		// header.SetModTime(time.Unix(info.ModTime().Unix(), 0))
		header.Modified = time.Unix(info.ModTime().Unix(), 0)
		//文件或者目录名
		header.Name = path

		//创建在zip内的文件或者目录
		writer, err := archive.CreateHeader(header)

		if err != nil {
			return err
		}

		//如果是目录，只需创建无需其他操作
		if info.IsDir() {
			return nil
		}

		//打开需要压缩的文件
		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		//将待压缩文件拷贝给zip内文件
		_, err = io.Copy(writer, file)

		return err

	})
}

/**
 * @description: 解压缩
 * @param {*} zipFile
 * @param {string} dest
 * @return {*}
 */
func DeCompress(zipFile, dest string) (err error) {
	//目标文件夹不存在则创建
	if _, err = os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dest, 0755)
		}
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		//    log.Println(file.Name)

		if file.FileInfo().IsDir() {

			err := os.MkdirAll(dest+"/"+file.Name, 0755)
			if err != nil {
				log.Println(err)
			}
			continue
		} else {
			res, err := getDir(dest + "/" + file.Name)
			if err != nil {
				return err
			}
			err = os.MkdirAll(res, 0755)
			if err != nil {
				return err
			}
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		filename := dest + "/" + file.Name
		//err = os.MkdirAll(getDir(filename), 0755)
		//if err != nil {
		//    return err
		//}

		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		//w.Close()
		//rc.Close()
	}
	return
}

func getDir(path string) (string, error) {
	res, err := subString(path, 0, strings.LastIndex(path, "/"))
	if err != nil {
		return res, errors.Wrapf(err, "getDir()")
	}
	return res, nil
}

func subString(str string, start, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", errors.New("subString start is wrong")
	}

	if end < start || end > length {
		return "", errors.New("subString end is wrong")
	}

	return string(rs[start:end]), nil
}
