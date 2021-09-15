/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-15 15:35:44
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 17:58:15
 * @FilePath: /goabao/abaodocker/model.go
 * @Description:
 */
package abaodocker

import "github.com/docker/docker/api/types"

type AuthConfigMap map[string]types.AuthConfig

// DockerImageBaseConfig 容器镜像基础信息配置
type DockerImageBaseConfig struct {
	RegistryServerURI string //容器镜像服务地址
	RegistryUserName  string //容器镜像服务用户名
	RegistryPassword  string //容器镜像服务密码
	ImageNamespace    string //容器镜像服务命名空间
	ImageName         string //容器镜像名称 ps: centos
	ImageVersion      string //容器镜像版本 ps: 7.8
}
