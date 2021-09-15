/*
 * @Author: lorock
 * @Github: https://github.com/lorock
 * @Date: 2021-09-14 15:29:20
 * @LastEditors: lorock
 * @LastEditTime: 2021-09-15 16:32:44
 * @FilePath: /goabao/abaodocker/image_push.go
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

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

/**
 * @description:
 * @param {DockerImageBaseConfig} dockerImageBaseConfig
 * @param {int64} timeOutSecond
 * @return {*}
 */
func ImagePush(dockerImageBaseConfig DockerImageBaseConfig, timeOutSecond int64) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSecond)*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrapf(err, "docker push error")
	}
	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: dockerImageBaseConfig.RegistryUserName,
		Password: dockerImageBaseConfig.RegistryPassword,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "docker push error")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	dockerImageUrl := []string{
		dockerImageBaseConfig.RegistryServerURI,
		dockerImageBaseConfig.ImageNamespace,
		dockerImageBaseConfig.ImageName,
	}

	dockerImageURI := strings.Join(dockerImageUrl, "/") + ":" + dockerImageBaseConfig.ImageVersion
	// docker images push
	iPush, iPushErr := cli.ImagePush(ctx, dockerImageURI, types.ImagePushOptions{RegistryAuth: authStr})

	if iPushErr != nil {
		return nil, errors.Wrapf(err, "docker push error")
	}

	defer iPush.Close()

	var messages []string

	rd := bufio.NewReader(iPush)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return messages, errors.Wrapf(err, "docker push image.bufio.NewReader()")
		}
		messages = append(messages, string(n))
	}
	return messages, nil

}
