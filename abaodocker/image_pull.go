/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-14 15:56:40
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 16:39:38
 * @FilePath: /goabao/abaodocker/image_pull.go
 * @Description:
 */
package abaodocker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

/**
 * @description:
 * @param {DockerImageBaseConfig} dockerImageBaseConfig
 * @param {int64} timeOutSecond
 * @return {*}
 */
func ImagePull(dockerImageBaseConfig DockerImageBaseConfig, timeOutSecond int64) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSecond)*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrapf(err, "docker pull cli error")
	}
	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: dockerImageBaseConfig.RegistryUserName,
		Password: dockerImageBaseConfig.RegistryPassword,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "docker pull encodedJSON error")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	dockerImageUrl := []string{
		dockerImageBaseConfig.RegistryServerURI,
		dockerImageBaseConfig.ImageNamespace,
		dockerImageBaseConfig.ImageName,
	}

	dockerImageURI := strings.Join(dockerImageUrl, "/") + ":" + dockerImageBaseConfig.ImageVersion

	// docker images pull
	iPull, err := cli.ImagePull(ctx, dockerImageURI, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return nil, errors.Wrapf(err, "docker pull out error")
	}

	defer iPull.Close()

	var messages []string

	rd := bufio.NewReader(iPull)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return messages, errors.Wrapf(err, "docker pull image.bufio.NewReader()")
		}
		messages = append(messages, string(n))
	}

	return messages, nil
}
