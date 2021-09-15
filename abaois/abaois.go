/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-15 15:15:35
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 15:15:36
 * @FilePath: /goabao/pkg/abaois/abaois.go
 * @Description:
 */
package abaois

import "os"

/**
 * @description: IsExist checks whether a file or directory exists,It returns false when the file or directory does not exist.
 * @param {string} f
 * @return {bool}
 */
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

/**
 * @description: IsFile checks whether the path is a file,it returns false when it's a directory or does not exist.
 * @param {*}
 * @return {*}
 */
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

/**
 * @description: IsDir checks whether the path is a dir,it returns false when it's a directory or does not exist.
 * @param {*}
 * @return {*}
 */
func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}
