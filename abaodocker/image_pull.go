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

func ImagePull(dockerRegistryServerURI, dockerUserName, dockerPassword, dockerNamespace, dockerName, dockerImageVersion string, timeOutSecond int64) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutSecond)*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrapf(err, "docker pull cli error")
	}
	defer cli.Close()

	authConfig := types.AuthConfig{
		Username: dockerUserName,
		Password: dockerPassword,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "docker pull encodedJSON error")
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	dockerUrl := []string{dockerRegistryServerURI, dockerNamespace, dockerName}

	// docker images pull
	iPull, err := cli.ImagePull(ctx, strings.Join(dockerUrl, "/")+":"+dockerImageVersion, types.ImagePullOptions{RegistryAuth: authStr})
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
