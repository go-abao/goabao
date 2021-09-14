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

func ImagePush(dockerRegistryServerURI, dockerUserName, dockerPassword, dockerNamespace, dockerName, dockerImageVersion string, timeOutSecond int64) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSecond)*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrapf(err, "docker push error")
	}
	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: dockerUserName,
		Password: dockerPassword,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "docker push error")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	dockerUrl := []string{dockerRegistryServerURI, dockerNamespace, dockerName}

	dockerURI := strings.Join(dockerUrl, "/") + ":" + dockerImageVersion
	// docker images push
	iPush, iPushErr := cli.ImagePush(ctx, dockerURI, types.ImagePushOptions{RegistryAuth: authStr})

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
