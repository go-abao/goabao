/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-14 15:19:42
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 18:05:49
 * @FilePath: /goabao/abaodocker/image_test.go
 * @Description:
 */
package abaodocker

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
)

var (
	dockerRegistryServerURI string = "ccr.ccs.tencentyun.com"
	dockerUserName          string = "username"
	dockerPassword          string = "password"
	dockerNamespace         string = "abao"
)

func TestImageBuild(t *testing.T) {

	dockerName := "go-test"
	dockerImageVersion := "1"

	srcDir := "test"

	authConfigs := make(AuthConfigMap)
	authConfigs[dockerRegistryServerURI] = types.AuthConfig{Username: dockerUserName, Password: dockerPassword}
	dockerImageBaseConfig := DockerImageBaseConfig{
		RegistryServerURI: dockerRegistryServerURI,
		RegistryUserName:  dockerUserName,
		RegistryPassword:  dockerPassword,
		ImageNamespace:    dockerNamespace,
		ImageName:         dockerName,
		ImageVersion:      dockerImageVersion,
	}
	msg, err := ImageBuild(srcDir, dockerImageBaseConfig, authConfigs, 300)
	if err != nil {
		t.Fatal(err)
	}

	for _, msgValue := range msg {
		t.Log(msgValue)
	}
	// fmt.Println(msg)

}

func TestImagePush(t *testing.T) {

	dockerName := "go-test"
	dockerImageVersion := "1"

	dockerImageBaseConfig := DockerImageBaseConfig{
		RegistryServerURI: dockerRegistryServerURI,
		RegistryUserName:  dockerUserName,
		RegistryPassword:  dockerPassword,
		ImageNamespace:    dockerNamespace,
		ImageName:         dockerName,
		ImageVersion:      dockerImageVersion,
	}

	out, err := ImagePush(dockerImageBaseConfig, 300)
	if err != nil {
		t.Fatalf("出错啦: %v", err)
	}
	if strings.Contains(strings.Join(out, ""), "error") {
		t.Fatalf("strings.ContainsAny error %v", strings.Join(out, "\n"))
	}
	for _, msgValue := range out {
		t.Log(msgValue)
	}
}

func TestImagePull(t *testing.T) {

	dockerName := "go-test"
	dockerImageVersion := "1"
	dockerImageBaseConfig := DockerImageBaseConfig{
		RegistryServerURI: dockerRegistryServerURI,
		RegistryUserName:  dockerUserName,
		RegistryPassword:  dockerPassword,
		ImageNamespace:    dockerNamespace,
		ImageName:         dockerName,
		ImageVersion:      dockerImageVersion,
	}
	out, err := ImagePull(dockerImageBaseConfig, 300)
	if err != nil {
		t.Fatalf("出错啦: %+v", err)
	}

	for _, msgValue := range out {
		t.Log(msgValue)
	}
}
